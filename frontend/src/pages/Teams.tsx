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
    const [newTeam, setNewTeam] = useState({ name: '', description: '' });

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
                                    <h3>Team Members ({members.length})</h3>
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
                                            </div>
                                        ))}
                                        {members.length === 0 && (
                                            <div className="empty-state">
                                                <p>No members in this team yet</p>
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

                {/* Create Modal */}
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
            </main>
        </div>
    );
};

export default Teams;
