import React, { useState, useEffect } from 'react';
import apiClient from '../api/client';
import Sidebar from '../components/Sidebar';
import './Board.css';

interface HoldReason {
    id: number;
    reason: string;
    created_at: string;
    resolved_at?: string;
    created_by_user?: {
        full_name: string;
    };
}

interface Issue {
    id: number;
    title: string;
    description: string;
    priority: string;
    status_id: number;
    created_by: number;
    is_on_hold?: boolean;
    hold_reasons?: HoldReason[];
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

interface Attachment {
    id: number;
    original_filename: string;
    file_size: number;
    mime_type: string;
    created_at: string;
    user?: { full_name: string };
}

const Board: React.FC = () => {
    const [teams, setTeams] = useState<Team[]>([]);
    const [selectedTeam, setSelectedTeam] = useState<number | null>(null);
    const [statuses, setStatuses] = useState<Status[]>([]);
    const [issues, setIssues] = useState<Issue[]>([]);
    const [loading, setLoading] = useState(true);
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [showHoldModal, setShowHoldModal] = useState(false);
    const [showDetailModal, setShowDetailModal] = useState(false);
    const [selectedIssue, setSelectedIssue] = useState<Issue | null>(null);
    const [holdReason, setHoldReason] = useState('');
    const [newIssue, setNewIssue] = useState({ title: '', description: '', priority: 'NORMAL' });
    const [attachments, setAttachments] = useState<Attachment[]>([]);
    const [uploading, setUploading] = useState(false);
    const fileInputRef = React.useRef<HTMLInputElement>(null);

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
                    apiClient.get(`/statuses`),
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

    const handleHoldIssue = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedIssue || !holdReason.trim()) return;

        try {
            await apiClient.post(`/issues/${selectedIssue.id}/hold`, { reason: holdReason });
            setIssues(issues.map(issue =>
                issue.id === selectedIssue.id ? { ...issue, is_on_hold: true } : issue
            ));
            setHoldReason('');
            setShowHoldModal(false);
            setSelectedIssue(null);
        } catch (error) {
            console.error('Failed to hold issue:', error);
        }
    };

    const handleResumeIssue = async (issueId: number) => {
        try {
            await apiClient.post(`/issues/${issueId}/resume`);
            setIssues(issues.map(issue =>
                issue.id === issueId ? { ...issue, is_on_hold: false } : issue
            ));
        } catch (error) {
            console.error('Failed to resume issue:', error);
        }
    };

    const openHoldModal = (issue: Issue) => {
        setSelectedIssue(issue);
        setShowHoldModal(true);
    };

    const openDetailModal = async (issue: Issue) => {
        try {
            const res = await apiClient.get(`/issues/${issue.id}`);
            setSelectedIssue(res.data);
            setShowDetailModal(true);
            // Fetch attachments
            try {
                const attachRes = await apiClient.get(`/issues/${issue.id}/attachments`);
                setAttachments(attachRes.data || []);
            } catch {
                setAttachments([]);
            }
        } catch (error) {
            console.error('Failed to fetch issue details:', error);
        }
    };

    const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
        if (!e.target.files || !e.target.files[0] || !selectedIssue) return;
        const file = e.target.files[0];

        if (file.size > 10 * 1024 * 1024) {
            alert('File too large. Maximum size is 10MB');
            return;
        }

        setUploading(true);
        const formData = new FormData();
        formData.append('file', file);

        try {
            const res = await apiClient.post(`/issues/${selectedIssue.id}/attachments`, formData, {
                headers: { 'Content-Type': 'multipart/form-data' }
            });
            setAttachments([res.data, ...attachments]);
        } catch (error: any) {
            alert(error.response?.data?.error || 'Failed to upload file');
        } finally {
            setUploading(false);
            if (fileInputRef.current) fileInputRef.current.value = '';
        }
    };

    const handleDownloadAttachment = async (attachment: Attachment) => {
        try {
            const res = await apiClient.get(`/attachments/${attachment.id}/download`);
            window.open(res.data.url, '_blank');
        } catch (error) {
            alert('Failed to download file');
        }
    };

    const handleDeleteAttachment = async (attachmentId: number) => {
        if (!confirm('Delete this attachment?')) return;
        try {
            await apiClient.delete(`/attachments/${attachmentId}`);
            setAttachments(attachments.filter(a => a.id !== attachmentId));
        } catch (error) {
            alert('Failed to delete attachment');
        }
    };

    const formatFileSize = (bytes: number) => {
        if (bytes < 1024) return bytes + ' B';
        if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
        return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
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
                        <p>Manage tasks and track progress</p>
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
                                    <div
                                        key={issue.id}
                                        className={`task-card ${issue.is_on_hold ? 'on-hold' : ''}`}
                                        onClick={() => openDetailModal(issue)}
                                    >
                                        {issue.is_on_hold && (
                                            <div className="hold-badge">⏸️ ON HOLD</div>
                                        )}
                                        <div className={`task-priority ${getPriorityClass(issue.priority)}`}>
                                            {issue.priority}
                                        </div>
                                        <h4>{issue.title}</h4>
                                        {issue.description && (
                                            <p className="task-desc">{issue.description.substring(0, 80)}...</p>
                                        )}
                                        <div className="task-actions" onClick={e => e.stopPropagation()}>
                                            <select
                                                value={issue.status_id}
                                                onChange={(e) => handleStatusChange(issue.id, Number(e.target.value))}
                                                className="status-select"
                                            >
                                                {statuses.map(s => (
                                                    <option key={s.id} value={s.id}>{s.name}</option>
                                                ))}
                                            </select>
                                            {issue.is_on_hold ? (
                                                <button
                                                    className="btn-action resume"
                                                    onClick={() => handleResumeIssue(issue.id)}
                                                    title="Resume task"
                                                >
                                                    ▶️
                                                </button>
                                            ) : (
                                                <button
                                                    className="btn-action hold"
                                                    onClick={() => openHoldModal(issue)}
                                                    title="Put on hold"
                                                >
                                                    ⏸️
                                                </button>
                                            )}
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

                {/* Hold Modal */}
                {showHoldModal && selectedIssue && (
                    <div className="modal-overlay" onClick={() => { setShowHoldModal(false); setHoldReason(''); }}>
                        <div className="modal" onClick={e => e.stopPropagation()}>
                            <h2>⏸️ Put Task On Hold</h2>
                            <p className="modal-subtitle">Task: {selectedIssue.title}</p>
                            <form onSubmit={handleHoldIssue}>
                                <div className="form-group">
                                    <label className="form-label">Reason for Hold *</label>
                                    <textarea
                                        className="form-input"
                                        rows={4}
                                        placeholder="Why is this task being put on hold? e.g., Waiting for client feedback, Blocked by another task, etc."
                                        value={holdReason}
                                        onChange={(e) => setHoldReason(e.target.value)}
                                        required
                                    />
                                </div>
                                <div className="modal-actions">
                                    <button type="button" className="btn btn-secondary" onClick={() => { setShowHoldModal(false); setHoldReason(''); }}>
                                        Cancel
                                    </button>
                                    <button type="submit" className="btn btn-warning">
                                        Put On Hold
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}

                {/* Detail Modal */}
                {showDetailModal && selectedIssue && (
                    <div className="modal-overlay" onClick={() => setShowDetailModal(false)}>
                        <div className="modal modal-lg" onClick={e => e.stopPropagation()}>
                            <div className="modal-header">
                                <h2>{selectedIssue.title}</h2>
                                <button className="btn-close" onClick={() => setShowDetailModal(false)}>✕</button>
                            </div>

                            <div className="detail-grid">
                                <div className="detail-item">
                                    <label>Priority</label>
                                    <span className={`badge ${getPriorityClass(selectedIssue.priority)}`}>
                                        {selectedIssue.priority}
                                    </span>
                                </div>
                                <div className="detail-item">
                                    <label>Status</label>
                                    <span>{statuses.find(s => s.id === selectedIssue.status_id)?.name || 'Unknown'}</span>
                                </div>
                                <div className="detail-item">
                                    <label>On Hold</label>
                                    <span>{selectedIssue.is_on_hold ? '⏸️ Yes' : '▶️ No'}</span>
                                </div>
                            </div>

                            {selectedIssue.description && (
                                <div className="detail-section">
                                    <h3>Description</h3>
                                    <p>{selectedIssue.description}</p>
                                </div>
                            )}

                            {selectedIssue.hold_reasons && selectedIssue.hold_reasons.length > 0 && (
                                <div className="detail-section">
                                    <h3>🚫 Hold History</h3>
                                    <div className="hold-history">
                                        {selectedIssue.hold_reasons.map((hr, idx) => (
                                            <div key={idx} className={`hold-item ${!hr.resolved_at ? 'active' : 'resolved'}`}>
                                                <div className="hold-reason">
                                                    <strong>{!hr.resolved_at ? '⏸️ Current Hold' : '✅ Resolved'}</strong>
                                                    <p>{hr.reason}</p>
                                                </div>
                                                <div className="hold-meta">
                                                    <span>By: {hr.created_by_user?.full_name || 'Unknown'}</span>
                                                    <span>{new Date(hr.created_at).toLocaleDateString()}</span>
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            )}

                            {/* Attachments Section */}
                            <div className="detail-section">
                                <div className="attachments-header">
                                    <h3>📎 Attachments</h3>
                                    <input
                                        type="file"
                                        ref={fileInputRef}
                                        onChange={handleFileUpload}
                                        style={{ display: 'none' }}
                                    />
                                    <button
                                        className="btn btn-secondary btn-sm"
                                        onClick={() => fileInputRef.current?.click()}
                                        disabled={uploading}
                                    >
                                        {uploading ? '⏳ Uploading...' : '+ Add File'}
                                    </button>
                                </div>
                                {attachments.length > 0 ? (
                                    <div className="attachments-list">
                                        {attachments.map(att => (
                                            <div key={att.id} className="attachment-item">
                                                <div className="attachment-info">
                                                    <span className="attachment-name" title={att.original_filename}>
                                                        📄 {att.original_filename}
                                                    </span>
                                                    <span className="attachment-meta">
                                                        {formatFileSize(att.file_size)} • {new Date(att.created_at).toLocaleDateString()}
                                                    </span>
                                                </div>
                                                <div className="attachment-actions">
                                                    <button
                                                        type="button"
                                                        className="btn btn-sm"
                                                        onClick={(e) => { e.stopPropagation(); handleDownloadAttachment(att); }}
                                                        title="Download"
                                                    >
                                                        ⬇️
                                                    </button>
                                                    <button
                                                        type="button"
                                                        className="btn btn-sm btn-danger"
                                                        onClick={(e) => { e.stopPropagation(); handleDeleteAttachment(att.id); }}
                                                        title="Delete"
                                                    >
                                                        🗑️
                                                    </button>
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                ) : (
                                    <p className="no-attachments">No attachments yet</p>
                                )}
                            </div>

                            <div className="modal-actions">
                                {selectedIssue.is_on_hold ? (
                                    <button
                                        className="btn btn-success"
                                        onClick={() => { handleResumeIssue(selectedIssue.id); setShowDetailModal(false); }}
                                    >
                                        ▶️ Resume Task
                                    </button>
                                ) : (
                                    <button
                                        className="btn btn-warning"
                                        onClick={() => { setShowDetailModal(false); openHoldModal(selectedIssue); }}
                                    >
                                        ⏸️ Put On Hold
                                    </button>
                                )}
                                <button className="btn btn-secondary" onClick={() => setShowDetailModal(false)}>
                                    Close
                                </button>
                            </div>
                        </div>
                    </div>
                )}
            </main>
        </div>
    );
};

export default Board;
