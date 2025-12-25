import React, { useState, useEffect } from 'react';
import apiClient from '../api/client';
import Sidebar from '../components/Sidebar';
import './Teams.css';

interface Team {
    id: number;
    name: string;
    description: string;
    parent_team_id?: number;
}

interface TeamMember {
    id: number;
    user_id: number;
    team_id: number;
    role: string;
    user?: {
        id: number;
        full_name: string;
        email: string;
    };
}

const Teams: React.FC = () => {
    const [teams, setTeams] = useState<Team[]>([]);
    const [selectedTeam, setSelectedTeam] = useState<Team | null>(null);
    const [members, setMembers] = useState<TeamMember[]>([]);
    const [loading, setLoading] = useState(true);
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [showAddMemberModal, setShowAddMemberModal] = useState(false);
    const [newTeam, setNewTeam] = useState({ name: '', description: '' });
    const [newMember, setNewMember] = useState({ email: '', role: 'member' });
    const [addMemberError, setAddMemberError] = useState('');
    const [addMemberSuccess, setAddMemberSuccess] = useState('');

    useEffect(() => {
        fetchTeams();
    }, []);

    const fetchTeams = async () => {
        try {
            const res = await apiClient.get('/teams');
            setTeams(res.data || []);
            if (res.data && res.data.length > 0) {
                await selectTeam(res.data[0]);
            }
        } catch (error) {
            console.error('Failed to fetch teams:', error);
        } finally {
            setLoading(false);
        }
    };

    const selectTeam = async (team: Team) => {
        setSelectedTeam(team);
        try {
            const res = await apiClient.get(`/teams/${team.id}/members`);
            setMembers(res.data || []);
        } catch (error) {
            console.error('Failed to fetch members:', error);
            setMembers([]);
        }
    };

    const handleCreateTeam = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const res = await apiClient.post('/teams', newTeam);
            setTeams([...teams, res.data]);
            setNewTeam({ name: '', description: '' });
            setShowCreateModal(false);
        } catch (error) {
            console.error('Failed to create team:', error);
        }
    };

    const handleAddMember = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedTeam) return;

        setAddMemberError('');
        setAddMemberSuccess('');

        try {
            await apiClient.post(`/teams/${selectedTeam.id}/members`, {
                email: newMember.email,
                role: newMember.role
            });

            // Refresh members
            const res = await apiClient.get(`/teams/${selectedTeam.id}/members`);
            setMembers(res.data || []);

            setAddMemberSuccess(`Successfully added ${newMember.email} to the team!`);
            setNewMember({ email: '', role: 'member' });

            // Auto close after success
            setTimeout(() => {
                setShowAddMemberModal(false);
                setAddMemberSuccess('');
            }, 1500);
        } catch (error: any) {
            setAddMemberError(error.response?.data?.error || 'Failed to add member. Make sure the email exists.');
        }
    };

    const handleRemoveMember = async (userId: number) => {
        if (!selectedTeam) return;
        if (!confirm('Are you sure you want to remove this member?')) return;

        try {
            await apiClient.delete(`/teams/${selectedTeam.id}/members/${userId}`);
            setMembers(members.filter(m => m.user_id !== userId));
        } catch (error) {
            console.error('Failed to remove member:', error);
        }
    };

    const getRoleBadgeClass = (role: string) => {
        switch (role) {
            case 'manager': return 'badge-primary';
            case 'assistant': return 'badge-warning';
            case 'member': return 'badge-success';
            default: return 'badge-secondary';
        }
    };

    if (loading) {
        return (
            <div className="page-layout">
                <Sidebar activeItem="teams" />
                <main className="main-content">
                    <div className="loading-screen">
                        <div className="spinner-lg"></div>
                        <p>Loading teams...</p>
                    </div>
                </main>
            </div>
        );
    }

    return (
        <div className="page-layout">
            <Sidebar activeItem="teams" />
            <main className="main-content">
                <header className="page-header">
                    <div className="header-left">
                        <h1>Teams</h1>
                        <p>Manage your teams and members</p>
                    </div>
                    <div className="header-actions">
                        <button className="btn btn-primary" onClick={() => setShowCreateModal(true)}>
                            + New Team
                        </button>
                    </div>
                </header>

                <div className="teams-layout">
                    {/* Teams List */}
                    <div className="teams-list">
                        <h3>All Teams</h3>
                        {teams.map(team => (
                            <div
                                key={team.id}
                                className={`team-card ${selectedTeam?.id === team.id ? 'active' : ''}`}
                                onClick={() => selectTeam(team)}
                            >
                                <div className="team-avatar">{team.name.charAt(0)}</div>
                                <div className="team-info">
                                    <h4>{team.name}</h4>
                                    <p>{team.description || 'No description'}</p>
                                </div>
                            </div>
                        ))}
                    </div>

                    {/* Team Details */}
                    <div className="team-details">
                        {selectedTeam ? (
                            <>
                                <div className="details-header">
                                    <div className="team-avatar-lg">{selectedTeam.name.charAt(0)}</div>
                                    <div>
                                        <h2>{selectedTeam.name}</h2>
                                        <p>{selectedTeam.description}</p>
                                    </div>
                                </div>

                                <div className="members-section">
                                    <div className="members-header">
                                        <h3>Team Members ({members.length})</h3>
                                        <button
                                            className="btn btn-primary"
                                            onClick={() => setShowAddMemberModal(true)}
                                        >
                                            + Add Member
                                        </button>
                                    </div>
                                    <div className="members-grid">
                                        {members.map(member => (
                                            <div key={member.id} className="member-card">
                                                <div className="member-avatar">
                                                    {member.user?.full_name?.charAt(0) || 'U'}
                                                </div>
                                                <div className="member-info">
                                                    <h4>{member.user?.full_name || 'Unknown'}</h4>
                                                    <p>{member.user?.email}</p>
                                                    <span className={`badge ${getRoleBadgeClass(member.role)}`}>
                                                        {member.role}
                                                    </span>
                                                </div>
                                                <button
                                                    className="btn-remove"
                                                    onClick={() => handleRemoveMember(member.user_id)}
                                                    title="Remove member"
                                                >
                                                    ✕
                                                </button>
                                            </div>
                                        ))}
                                        {members.length === 0 && (
                                            <div className="empty-state">
                                                <p>No members in this team yet</p>
                                                <button
                                                    className="btn btn-primary"
                                                    onClick={() => setShowAddMemberModal(true)}
                                                >
                                                    Add First Member
                                                </button>
                                            </div>
                                        )}
                                    </div>
                                </div>
                            </>
                        ) : (
                            <div className="empty-state">
                                <p>Select a team to view details</p>
                            </div>
                        )}
                    </div>
                </div>

                {/* Create Team Modal */}
                {showCreateModal && (
                    <div className="modal-overlay" onClick={() => setShowCreateModal(false)}>
                        <div className="modal" onClick={e => e.stopPropagation()}>
                            <h2>Create New Team</h2>
                            <form onSubmit={handleCreateTeam}>
                                <div className="form-group">
                                    <label className="form-label">Team Name</label>
                                    <input
                                        type="text"
                                        className="form-input"
                                        value={newTeam.name}
                                        onChange={(e) => setNewTeam({ ...newTeam, name: e.target.value })}
                                        required
                                    />
                                </div>
                                <div className="form-group">
                                    <label className="form-label">Description</label>
                                    <textarea
                                        className="form-input"
                                        rows={3}
                                        value={newTeam.description}
                                        onChange={(e) => setNewTeam({ ...newTeam, description: e.target.value })}
                                    />
                                </div>
                                <div className="modal-actions">
                                    <button type="button" className="btn btn-secondary" onClick={() => setShowCreateModal(false)}>
                                        Cancel
                                    </button>
                                    <button type="submit" className="btn btn-primary">
                                        Create Team
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}

                {/* Add Member Modal */}
                {showAddMemberModal && (
                    <div className="modal-overlay" onClick={() => { setShowAddMemberModal(false); setAddMemberError(''); }}>
                        <div className="modal" onClick={e => e.stopPropagation()}>
                            <h2>Add Team Member</h2>
                            <p className="modal-subtitle">Add a user to {selectedTeam?.name}</p>

                            {addMemberError && (
                                <div className="alert alert-error">
                                    <span>⚠️</span> {addMemberError}
                                </div>
                            )}

                            {addMemberSuccess && (
                                <div className="alert alert-success">
                                    <span>✅</span> {addMemberSuccess}
                                </div>
                            )}

                            <form onSubmit={handleAddMember}>
                                <div className="form-group">
                                    <label className="form-label">User Email</label>
                                    <input
                                        type="email"
                                        className="form-input"
                                        placeholder="Enter user email"
                                        value={newMember.email}
                                        onChange={(e) => setNewMember({ ...newMember, email: e.target.value })}
                                        required
                                    />
                                    <small className="form-hint">The user must already be registered in the system</small>
                                </div>
                                <div className="form-group">
                                    <label className="form-label">Role</label>
                                    <select
                                        className="form-input"
                                        value={newMember.role}
                                        onChange={(e) => setNewMember({ ...newMember, role: e.target.value })}
                                    >
                                        <option value="member">Member</option>
                                        <option value="assistant">Assistant</option>
                                        <option value="manager">Manager</option>
                                        <option value="stakeholder">Stakeholder</option>
                                    </select>
                                </div>
                                <div className="modal-actions">
                                    <button type="button" className="btn btn-secondary" onClick={() => { setShowAddMemberModal(false); setAddMemberError(''); }}>
                                        Cancel
                                    </button>
                                    <button type="submit" className="btn btn-primary">
                                        Add Member
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}
            </main>
        </div>
    );
};

export default Teams;
