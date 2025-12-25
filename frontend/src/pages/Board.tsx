import React, { useState, useEffect } from 'react';
import apiClient from '../api/client';
import Sidebar from '../components/Sidebar';
import './Board.css';

interface Issue {
    id: number;
    title: string;
    description: string;
    priority: string;
    status_id: number;
    created_by: number;
}

interface Status {
    id: number;
    name: string;
    color: string;
    position: number;
    is_final: boolean;
}

interface Team {
    id: number;
    name: string;
}

const Board: React.FC = () => {
    const [teams, setTeams] = useState<Team[]>([]);
    const [selectedTeam, setSelectedTeam] = useState<number | null>(null);
    const [statuses, setStatuses] = useState<Status[]>([]);
    const [issues, setIssues] = useState<Issue[]>([]);
    const [loading, setLoading] = useState(true);
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [newIssue, setNewIssue] = useState({ title: '', description: '', priority: 'NORMAL' });

    useEffect(() => {
        const fetchTeams = async () => {
            try {
                const res = await apiClient.get('/teams');
                setTeams(res.data || []);
                if (res.data && res.data.length > 0) {
                    setSelectedTeam(res.data[0].id);
                }
            } catch (error) {
                console.error('Failed to fetch teams:', error);
            }
        };
        fetchTeams();
    }, []);

    useEffect(() => {
        if (!selectedTeam) return;

        const fetchBoardData = async () => {
            setLoading(true);
            try {
                const [statusRes, issuesRes] = await Promise.all([
                    apiClient.get(`/statuses/team/${selectedTeam}`),
                    apiClient.get(`/issues?team_id=${selectedTeam}`)
                ]);
                setStatuses((statusRes.data || []).sort((a: Status, b: Status) => a.position - b.position));
                setIssues(issuesRes.data || []);
            } catch (error) {
                console.error('Failed to fetch board data:', error);
            } finally {
                setLoading(false);
            }
        };
        fetchBoardData();
    }, [selectedTeam]);

    const handleStatusChange = async (issueId: number, newStatusId: number) => {
        try {
            await apiClient.post(`/issues/${issueId}/status`, { status_id: newStatusId });
            setIssues(issues.map(issue =>
                issue.id === issueId ? { ...issue, status_id: newStatusId } : issue
            ));
        } catch (error) {
            console.error('Failed to update status:', error);
        }
    };

    const handleCreateIssue = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedTeam || !statuses.length) return;

        try {
            const res = await apiClient.post('/issues', {
                team_id: selectedTeam,
                status_id: statuses[0].id,
                title: newIssue.title,
                description: newIssue.description,
                priority: newIssue.priority
            });
            setIssues([...issues, res.data]);
            setNewIssue({ title: '', description: '', priority: 'NORMAL' });
            setShowCreateModal(false);
        } catch (error) {
            console.error('Failed to create issue:', error);
        }
    };

    const getPriorityClass = (priority: string) => {
        switch (priority) {
            case 'URGENT': return 'priority-urgent';
            case 'HIGH': return 'priority-high';
            case 'NORMAL': return 'priority-normal';
            default: return 'priority-low';
        }
    };

    const getIssuesByStatus = (statusId: number) => {
        return issues.filter(issue => issue.status_id === statusId);
    };

    if (loading && !statuses.length) {
        return (
            <div className="page-layout">
                <Sidebar activeItem="board" />
                <main className="main-content">
                    <div className="loading-screen">
                        <div className="spinner-lg"></div>
                        <p>Loading board...</p>
                    </div>
                </main>
            </div>
        );
    }

    return (
        <div className="page-layout">
            <Sidebar activeItem="board" />
            <main className="main-content">
                <header className="page-header">
                    <div className="header-left">
                        <h1>Kanban Board</h1>
                        <p>Drag tasks between columns to update status</p>
                    </div>
                    <div className="header-actions">
                        <select
                            className="team-selector"
                            value={selectedTeam || ''}
                            onChange={(e) => setSelectedTeam(Number(e.target.value))}
                        >
                            {teams.map(team => (
                                <option key={team.id} value={team.id}>{team.name}</option>
                            ))}
                        </select>
                        <button className="btn btn-primary" onClick={() => setShowCreateModal(true)}>
                            + New Task
                        </button>
                    </div>
                </header>

                <div className="board-container">
                    {statuses.map(status => (
                        <div key={status.id} className="board-column">
                            <div className="column-header" style={{ borderTopColor: status.color }}>
                                <h3>{status.name}</h3>
                                <span className="task-count">{getIssuesByStatus(status.id).length}</span>
                            </div>
                            <div className="column-content">
                                {getIssuesByStatus(status.id).map(issue => (
                                    <div key={issue.id} className="task-card">
                                        <div className={`task-priority ${getPriorityClass(issue.priority)}`}>
                                            {issue.priority}
                                        </div>
                                        <h4>{issue.title}</h4>
                                        {issue.description && (
                                            <p className="task-desc">{issue.description.substring(0, 80)}...</p>
                                        )}
                                        <div className="task-actions">
                                            <select
                                                value={issue.status_id}
                                                onChange={(e) => handleStatusChange(issue.id, Number(e.target.value))}
                                                className="status-select"
                                            >
                                                {statuses.map(s => (
                                                    <option key={s.id} value={s.id}>{s.name}</option>
                                                ))}
                                            </select>
                                        </div>
                                    </div>
                                ))}
                                {getIssuesByStatus(status.id).length === 0 && (
                                    <div className="empty-column">No tasks</div>
                                )}
                            </div>
                        </div>
                    ))}
                </div>

                {/* Create Modal */}
                {showCreateModal && (
                    <div className="modal-overlay" onClick={() => setShowCreateModal(false)}>
                        <div className="modal" onClick={e => e.stopPropagation()}>
                            <h2>Create New Task</h2>
                            <form onSubmit={handleCreateIssue}>
                                <div className="form-group">
                                    <label className="form-label">Title</label>
                                    <input
                                        type="text"
                                        className="form-input"
                                        value={newIssue.title}
                                        onChange={(e) => setNewIssue({ ...newIssue, title: e.target.value })}
                                        required
                                    />
                                </div>
                                <div className="form-group">
                                    <label className="form-label">Description</label>
                                    <textarea
                                        className="form-input"
                                        rows={4}
                                        value={newIssue.description}
                                        onChange={(e) => setNewIssue({ ...newIssue, description: e.target.value })}
                                    />
                                </div>
                                <div className="form-group">
                                    <label className="form-label">Priority</label>
                                    <select
                                        className="form-input"
                                        value={newIssue.priority}
                                        onChange={(e) => setNewIssue({ ...newIssue, priority: e.target.value })}
                                    >
                                        <option value="LOW">Low</option>
                                        <option value="NORMAL">Normal</option>
                                        <option value="HIGH">High</option>
                                        <option value="URGENT">Urgent</option>
                                    </select>
                                </div>
                                <div className="modal-actions">
                                    <button type="button" className="btn btn-secondary" onClick={() => setShowCreateModal(false)}>
                                        Cancel
                                    </button>
                                    <button type="submit" className="btn btn-primary">
                                        Create Task
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

export default Board;
