import React, { useState, useEffect } from 'react';
import apiClient from '../api/client';
import Sidebar from '../components/Sidebar';
import './Meetings.css';

interface Attendee {
    id: number;
    user_id: number;
    status: string;
    user?: { full_name: string };
}

interface Meeting {
    id: number;
    team_id: number;
    title: string;
    description: string;
    meeting_date: string;
    start_time: string;
    end_time: string;
    location: string;
    is_recurring: boolean;
    recurring_pattern?: string;
    creator?: { full_name: string };
    attendees?: Attendee[];
}

interface Team {
    id: number;
    name: string;
}

interface TeamMember {
    id: number;
    user_id: number;
    user?: { id: number; full_name: string; email: string };
    full_name?: string; // fallback
}

const Meetings: React.FC = () => {
    const [meetings, setMeetings] = useState<Meeting[]>([]);
    const [teams, setTeams] = useState<Team[]>([]);
    const [teamMembers, setTeamMembers] = useState<TeamMember[]>([]);
    const [selectedTeam, setSelectedTeam] = useState<number | null>(null);
    const [loading, setLoading] = useState(true);
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [selectedMeeting, setSelectedMeeting] = useState<Meeting | null>(null);
    const [newMeeting, setNewMeeting] = useState({
        title: '',
        description: '',
        meeting_date: '',
        start_time: '',
        end_time: '',
        location: '',
        is_recurring: false,
        recurring_pattern: '',
        attendee_ids: [] as number[]
    });

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

        const fetchData = async () => {
            setLoading(true);
            try {
                const [meetingsRes, membersRes] = await Promise.all([
                    apiClient.get(`/meetings?team_id=${selectedTeam}`),
                    apiClient.get(`/teams/${selectedTeam}/members`)
                ]);
                setMeetings(meetingsRes.data || []);
                // Map members to extract user info
                const members = (membersRes.data || []).map((m: TeamMember) => ({
                    ...m,
                    // Use user.id if available, else user_id
                    id: m.user?.id || m.user_id || m.id,
                    full_name: m.user?.full_name || m.full_name || 'Unknown'
                }));
                setTeamMembers(members);
            } catch (error) {
                console.error('Failed to fetch data:', error);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [selectedTeam]);

    const handleCreateMeeting = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedTeam) return;

        try {
            const res = await apiClient.post('/meetings', {
                team_id: selectedTeam,
                ...newMeeting
            });
            setMeetings([res.data, ...meetings]);
            setShowCreateModal(false);
            setNewMeeting({
                title: '',
                description: '',
                meeting_date: '',
                start_time: '',
                end_time: '',
                location: '',
                is_recurring: false,
                recurring_pattern: '',
                attendee_ids: []
            });
        } catch (error) {
            alert('Failed to create meeting');
        }
    };

    const handleDeleteMeeting = async (id: number) => {
        if (!confirm('Delete this meeting?')) return;
        try {
            await apiClient.delete(`/meetings/${id}`);
            setMeetings(meetings.filter(m => m.id !== id));
            setSelectedMeeting(null);
        } catch (error) {
            alert('Failed to delete meeting');
        }
    };

    const formatDate = (dateStr: string) => {
        return new Date(dateStr).toLocaleDateString('en-US', {
            weekday: 'short',
            month: 'short',
            day: 'numeric'
        });
    };

    const formatTime = (time: string) => {
        const [hours, minutes] = time.split(':');
        const h = parseInt(hours);
        const ampm = h >= 12 ? 'PM' : 'AM';
        const h12 = h % 12 || 12;
        return `${h12}:${minutes} ${ampm}`;
    };

    const getStatusBadge = (status: string) => {
        switch (status) {
            case 'ACCEPTED': return 'badge-success';
            case 'DECLINED': return 'badge-danger';
            default: return 'badge-secondary';
        }
    };

    // Group meetings by date
    const groupedMeetings = meetings.reduce((acc, meeting) => {
        const date = meeting.meeting_date.split('T')[0];
        if (!acc[date]) acc[date] = [];
        acc[date].push(meeting);
        return acc;
    }, {} as Record<string, Meeting[]>);

    return (
        <div className="page-layout">
            <Sidebar activeItem="meetings" />
            <main className="main-content">
                <header className="page-header">
                    <div className="header-left">
                        <h1>📅 Meetings</h1>
                        <p>Schedule and manage team meetings</p>
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
                        <button
                            className="btn btn-primary"
                            onClick={() => setShowCreateModal(true)}
                        >
                            + New Meeting
                        </button>
                    </div>
                </header>

                {loading ? (
                    <div className="loading-screen">
                        <div className="spinner-lg"></div>
                    </div>
                ) : meetings.length === 0 ? (
                    <div className="empty-state">
                        <span className="empty-icon">📅</span>
                        <h3>No meetings scheduled</h3>
                        <p>Create your first meeting to get started</p>
                    </div>
                ) : (
                    <div className="meetings-timeline">
                        {Object.entries(groupedMeetings).sort().map(([date, dayMeetings]) => (
                            <div key={date} className="meeting-day">
                                <div className="day-header">
                                    <span className="day-date">{formatDate(date)}</span>
                                    <span className="day-count">{dayMeetings.length} meeting(s)</span>
                                </div>
                                <div className="day-meetings">
                                    {dayMeetings.map(meeting => (
                                        <div
                                            key={meeting.id}
                                            className="meeting-card"
                                            onClick={() => setSelectedMeeting(meeting)}
                                        >
                                            <div className="meeting-time">
                                                <span className="time-start">{formatTime(meeting.start_time)}</span>
                                                <span className="time-sep">-</span>
                                                <span className="time-end">{formatTime(meeting.end_time)}</span>
                                            </div>
                                            <div className="meeting-info">
                                                <h3>{meeting.title}</h3>
                                                {meeting.location && (
                                                    <p className="meeting-location">📍 {meeting.location}</p>
                                                )}
                                                {meeting.attendees && meeting.attendees.length > 0 && (
                                                    <div className="meeting-attendees">
                                                        👥 {meeting.attendees.length} attendee(s)
                                                    </div>
                                                )}
                                            </div>
                                            {meeting.is_recurring && (
                                                <span className="recurring-badge">🔁 {meeting.recurring_pattern}</span>
                                            )}
                                        </div>
                                    ))}
                                </div>
                            </div>
                        ))}
                    </div>
                )}

                {/* Create Meeting Modal */}
                {showCreateModal && (
                    <div className="modal-overlay" onClick={() => setShowCreateModal(false)}>
                        <div className="modal" onClick={e => e.stopPropagation()}>
                            <div className="modal-header">
                                <h2>📅 New Meeting</h2>
                                <button className="btn-close" onClick={() => setShowCreateModal(false)}>✕</button>
                            </div>
                            <form onSubmit={handleCreateMeeting}>
                                <div className="form-group">
                                    <label className="form-label">Title *</label>
                                    <input
                                        type="text"
                                        className="form-input"
                                        placeholder="Meeting title"
                                        value={newMeeting.title}
                                        onChange={e => setNewMeeting({ ...newMeeting, title: e.target.value })}
                                        required
                                    />
                                </div>
                                <div className="form-group">
                                    <label className="form-label">Description</label>
                                    <textarea
                                        className="form-input"
                                        rows={3}
                                        placeholder="Meeting description or agenda"
                                        value={newMeeting.description}
                                        onChange={e => setNewMeeting({ ...newMeeting, description: e.target.value })}
                                    />
                                </div>
                                <div className="form-row">
                                    <div className="form-group">
                                        <label className="form-label">Date *</label>
                                        <input
                                            type="date"
                                            className="form-input"
                                            value={newMeeting.meeting_date}
                                            onChange={e => setNewMeeting({ ...newMeeting, meeting_date: e.target.value })}
                                            required
                                        />
                                    </div>
                                    <div className="form-group">
                                        <label className="form-label">Start Time *</label>
                                        <input
                                            type="time"
                                            className="form-input"
                                            value={newMeeting.start_time}
                                            onChange={e => setNewMeeting({ ...newMeeting, start_time: e.target.value })}
                                            required
                                        />
                                    </div>
                                    <div className="form-group">
                                        <label className="form-label">End Time *</label>
                                        <input
                                            type="time"
                                            className="form-input"
                                            value={newMeeting.end_time}
                                            onChange={e => setNewMeeting({ ...newMeeting, end_time: e.target.value })}
                                            required
                                        />
                                    </div>
                                </div>
                                <div className="form-group">
                                    <label className="form-label">Location / Link</label>
                                    <input
                                        type="text"
                                        className="form-input"
                                        placeholder="Room name or video call link"
                                        value={newMeeting.location}
                                        onChange={e => setNewMeeting({ ...newMeeting, location: e.target.value })}
                                    />
                                </div>
                                <div className="form-group">
                                    <label className="form-label">Attendees</label>
                                    <div className="attendee-checkboxes">
                                        {teamMembers.map(member => (
                                            <label key={member.id} className="checkbox-label">
                                                <input
                                                    type="checkbox"
                                                    checked={newMeeting.attendee_ids.includes(member.id)}
                                                    onChange={e => {
                                                        if (e.target.checked) {
                                                            setNewMeeting({
                                                                ...newMeeting,
                                                                attendee_ids: [...newMeeting.attendee_ids, member.id]
                                                            });
                                                        } else {
                                                            setNewMeeting({
                                                                ...newMeeting,
                                                                attendee_ids: newMeeting.attendee_ids.filter(id => id !== member.id)
                                                            });
                                                        }
                                                    }}
                                                />
                                                {member.full_name}
                                            </label>
                                        ))}
                                    </div>
                                </div>
                                <div className="form-group">
                                    <label className="checkbox-label">
                                        <input
                                            type="checkbox"
                                            checked={newMeeting.is_recurring}
                                            onChange={e => setNewMeeting({ ...newMeeting, is_recurring: e.target.checked })}
                                        />
                                        Recurring meeting
                                    </label>
                                    {newMeeting.is_recurring && (
                                        <select
                                            className="form-input"
                                            value={newMeeting.recurring_pattern}
                                            onChange={e => setNewMeeting({ ...newMeeting, recurring_pattern: e.target.value })}
                                            style={{ marginTop: '0.5rem' }}
                                        >
                                            <option value="">Select pattern</option>
                                            <option value="DAILY">Daily</option>
                                            <option value="WEEKLY">Weekly</option>
                                            <option value="MONTHLY">Monthly</option>
                                        </select>
                                    )}
                                </div>
                                <div className="modal-actions">
                                    <button type="button" className="btn btn-secondary" onClick={() => setShowCreateModal(false)}>
                                        Cancel
                                    </button>
                                    <button type="submit" className="btn btn-primary">
                                        Create Meeting
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                )}

                {/* Meeting Detail Modal */}
                {selectedMeeting && (
                    <div className="modal-overlay" onClick={() => setSelectedMeeting(null)}>
                        <div className="modal" onClick={e => e.stopPropagation()}>
                            <div className="modal-header">
                                <h2>{selectedMeeting.title}</h2>
                                <button className="btn-close" onClick={() => setSelectedMeeting(null)}>✕</button>
                            </div>

                            <div className="meeting-details">
                                <div className="detail-row">
                                    <span className="detail-label">📅 Date</span>
                                    <span>{formatDate(selectedMeeting.meeting_date)}</span>
                                </div>
                                <div className="detail-row">
                                    <span className="detail-label">🕐 Time</span>
                                    <span>{formatTime(selectedMeeting.start_time)} - {formatTime(selectedMeeting.end_time)}</span>
                                </div>
                                {selectedMeeting.location && (
                                    <div className="detail-row">
                                        <span className="detail-label">📍 Location</span>
                                        <span>{selectedMeeting.location}</span>
                                    </div>
                                )}
                                {selectedMeeting.description && (
                                    <div className="detail-row description">
                                        <span className="detail-label">📝 Description</span>
                                        <p>{selectedMeeting.description}</p>
                                    </div>
                                )}
                                {selectedMeeting.creator && (
                                    <div className="detail-row">
                                        <span className="detail-label">👤 Organizer</span>
                                        <span>{selectedMeeting.creator.full_name}</span>
                                    </div>
                                )}

                                {selectedMeeting.attendees && selectedMeeting.attendees.length > 0 && (
                                    <div className="attendees-section">
                                        <h4>👥 Attendees ({selectedMeeting.attendees.length})</h4>
                                        <div className="attendees-list">
                                            {selectedMeeting.attendees.map(att => (
                                                <div key={att.id} className="attendee-item">
                                                    <span>{att.user?.full_name}</span>
                                                    <span className={`badge ${getStatusBadge(att.status)}`}>
                                                        {att.status}
                                                    </span>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                )}
                            </div>

                            <div className="modal-actions">
                                <button
                                    type="button"
                                    className="btn btn-danger"
                                    onClick={() => handleDeleteMeeting(selectedMeeting.id)}
                                >
                                    Delete Meeting
                                </button>
                            </div>
                        </div>
                    </div>
                )}
            </main>
        </div>
    );
};

export default Meetings;
