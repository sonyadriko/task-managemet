import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { Link } from 'react-router-dom';
import apiClient from '../api/client';
import Sidebar from '../components/Sidebar';
import { usePermissions } from '../contexts/PermissionContext';
import './Dashboard.css';

interface Team {
    id: number;
    name: string;
    description: string;
}

interface Issue {
    id: number;
    title: string;
    priority: string;
    status?: {
        name: string;
        color: string;
    };
}

interface StatusCount {
    status_id: number;
    status_name: string;
    color: string;
    count: number;
}

interface PriorityCount {
    priority: string;
    count: number;
}

interface DailyCount {
    date: string;
    count: number;
}

interface TeamStat {
    team_id: number;
    team_name: string;
    task_count: number;
    member_count: number;
}

interface Analytics {
    total_tasks: number;
    completed_tasks: number;
    in_progress_tasks: number;
    on_hold_tasks: number;
    overdue_tasks: number;
    tasks_by_status: StatusCount[];
    tasks_by_priority: PriorityCount[];
    weekly_activity: DailyCount[];
    team_stats: TeamStat[];
}

const Dashboard: React.FC = () => {
    const { user } = useAuth();
    const { permissions } = usePermissions();
    const [teams, setTeams] = useState<Team[]>([]);
    const [issues, setIssues] = useState<Issue[]>([]);
    const [analytics, setAnalytics] = useState<Analytics | null>(null);
    const [loading, setLoading] = useState(true);
    const [orgName, setOrgName] = useState('');

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [teamsRes, analyticsRes] = await Promise.all([
                    apiClient.get('/teams'),
                    apiClient.get('/analytics/dashboard')
                ]);
                setTeams(teamsRes.data || []);
                setAnalytics(analyticsRes.data);

                if (teamsRes.data && teamsRes.data.length > 0) {
                    const issuesRes = await apiClient.get(`/issues?team_id=${teamsRes.data[0].id}`);
                    setIssues(issuesRes.data || []);
                }

                if (user?.organization_id) {
                    const orgRes = await apiClient.get(`/organizations/${user.organization_id}`);
                    setOrgName(orgRes.data?.name || 'Unknown');
                }
            } catch (error) {
                console.error('Failed to fetch data:', error);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [user?.organization_id]);

    const getPriorityClass = (priority: string) => {
        switch (priority) {
            case 'URGENT': return 'badge-danger';
            case 'HIGH': return 'badge-warning';
            case 'NORMAL': return 'badge-primary';
            default: return 'badge-secondary';
        }
    };

    const getPriorityColor = (priority: string) => {
        switch (priority) {
            case 'URGENT': return '#ef4444';
            case 'HIGH': return '#f59e0b';
            case 'NORMAL': return '#4877BA';
            default: return '#6b7280';
        }
    };

    const getMaxWeeklyCount = () => {
        if (!analytics?.weekly_activity) return 1;
        return Math.max(...analytics.weekly_activity.map(d => d.count), 1);
    };

    const completionRate = analytics
        ? Math.round((analytics.completed_tasks / Math.max(analytics.total_tasks, 1)) * 100)
        : 0;

    if (loading) {
        return (
            <div className="page-layout">
                <Sidebar activeItem="dashboard" />
                <main className="main-content">
                    <div className="loading-screen">
                        <div className="spinner-lg"></div>
                        <p>Loading your workspace...</p>
                    </div>
                </main>
            </div>
        );
    }

    return (
        <div className="page-layout">
            <Sidebar activeItem="dashboard" />
            <main className="main-content">
                {/* Header */}
                <header className="page-header">
                    <div className="header-left">
                        <h1>Dashboard</h1>
                        <p>Welcome back, {user?.full_name}! 👋</p>
                    </div>
                    <div className="header-right">
                        <div className="user-menu">
                            <div className="user-avatar">
                                {user?.full_name?.charAt(0) || 'U'}
                            </div>
                            <div className="user-info">
                                <span className="user-name">{user?.full_name}</span>
                                <span className="user-email">{user?.email}</span>
                            </div>
                        </div>
                    </div>
                </header>

                {/* Stats Cards */}
                <div className="stats-grid animate-fadeIn">
                    <div className="stat-card">
                        <div className="stat-icon purple">📋</div>
                        <div className="stat-content">
                            <h3>{analytics?.total_tasks || 0}</h3>
                            <p>Total Tasks</p>
                        </div>
                    </div>
                    <div className="stat-card">
                        <div className="stat-icon blue">🔄</div>
                        <div className="stat-content">
                            <h3>{analytics?.in_progress_tasks || 0}</h3>
                            <p>In Progress</p>
                        </div>
                    </div>
                    <div className="stat-card">
                        <div className="stat-icon green">✅</div>
                        <div className="stat-content">
                            <h3>{analytics?.completed_tasks || 0}</h3>
                            <p>Completed</p>
                        </div>
                    </div>
                    <div className="stat-card">
                        <div className="stat-icon orange">⏸️</div>
                        <div className="stat-content">
                            <h3>{analytics?.on_hold_tasks || 0}</h3>
                            <p>On Hold</p>
                        </div>
                    </div>
                    {analytics && analytics.overdue_tasks > 0 && (
                        <div className="stat-card warning">
                            <div className="stat-icon red">⚠️</div>
                            <div className="stat-content">
                                <h3>{analytics.overdue_tasks}</h3>
                                <p>Overdue</p>
                            </div>
                        </div>
                    )}
                </div>

                {/* Analytics Charts */}
                <div className="analytics-grid animate-fadeIn">
                    {/* Completion Rate */}
                    <div className="analytics-card">
                        <h3>Completion Rate</h3>
                        <div className="completion-ring">
                            <svg viewBox="0 0 36 36">
                                <path
                                    d="M18 2.0845
                                    a 15.9155 15.9155 0 0 1 0 31.831
                                    a 15.9155 15.9155 0 0 1 0 -31.831"
                                    fill="none"
                                    stroke="var(--color-bg-tertiary)"
                                    strokeWidth="3"
                                />
                                <path
                                    d="M18 2.0845
                                    a 15.9155 15.9155 0 0 1 0 31.831
                                    a 15.9155 15.9155 0 0 1 0 -31.831"
                                    fill="none"
                                    stroke="var(--color-success)"
                                    strokeWidth="3"
                                    strokeDasharray={`${completionRate}, 100`}
                                    strokeLinecap="round"
                                />
                            </svg>
                            <div className="completion-text">
                                <span className="percentage">{completionRate}%</span>
                                <span className="label">Complete</span>
                            </div>
                        </div>
                    </div>

                    {/* Tasks by Status */}
                    <div className="analytics-card">
                        <h3>Tasks by Status</h3>
                        <div className="status-bars">
                            {analytics?.tasks_by_status?.map(status => (
                                <div key={status.status_id} className="status-bar-item">
                                    <div className="status-bar-label">
                                        <span className="status-dot" style={{ backgroundColor: status.color }}></span>
                                        <span>{status.status_name}</span>
                                        <span className="count">{status.count}</span>
                                    </div>
                                    <div className="status-bar-track">
                                        <div
                                            className="status-bar-fill"
                                            style={{
                                                width: `${(status.count / Math.max(analytics.total_tasks, 1)) * 100}%`,
                                                backgroundColor: status.color
                                            }}
                                        ></div>
                                    </div>
                                </div>
                            ))}
                            {(!analytics?.tasks_by_status || analytics.tasks_by_status.length === 0) && (
                                <p className="no-data">No status data available</p>
                            )}
                        </div>
                    </div>

                    {/* Tasks by Priority */}
                    <div className="analytics-card">
                        <h3>Tasks by Priority</h3>
                        <div className="priority-chart">
                            {analytics?.tasks_by_priority?.map(p => (
                                <div key={p.priority} className="priority-item">
                                    <div
                                        className="priority-bar"
                                        style={{
                                            height: `${Math.max((p.count / Math.max(analytics.total_tasks, 1)) * 100, 10)}%`,
                                            backgroundColor: getPriorityColor(p.priority)
                                        }}
                                    >
                                        <span className="priority-count">{p.count}</span>
                                    </div>
                                    <span className="priority-label">{p.priority}</span>
                                </div>
                            ))}
                            {(!analytics?.tasks_by_priority || analytics.tasks_by_priority.length === 0) && (
                                <p className="no-data">No priority data available</p>
                            )}
                        </div>
                    </div>

                    {/* Weekly Activity */}
                    <div className="analytics-card wide">
                        <h3>Weekly Activity</h3>
                        <div className="activity-chart">
                            {analytics?.weekly_activity?.map((day, i) => (
                                <div key={i} className="activity-bar-item">
                                    <div
                                        className="activity-bar"
                                        style={{ height: `${Math.max((day.count / getMaxWeeklyCount()) * 100, 5)}%` }}
                                    >
                                        {day.count > 0 && <span className="activity-count">{day.count}</span>}
                                    </div>
                                    <span className="activity-day">{day.date}</span>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>

                {/* Team Stats & Recent Tasks */}
                <div className="dashboard-grid">
                    {/* Team Stats */}
                    {analytics?.team_stats && analytics.team_stats.length > 0 && (
                        <section className="dashboard-section animate-fadeIn">
                            <div className="section-header">
                                <h2>Team Overview</h2>
                            </div>
                            <div className="team-stats-grid">
                                {analytics.team_stats.map(team => (
                                    <div key={team.team_id} className="team-stat-card">
                                        <div className="team-stat-header">
                                            <span className="team-avatar">{team.team_name.charAt(0)}</span>
                                            <h4>{team.team_name}</h4>
                                        </div>
                                        <div className="team-stat-numbers">
                                            <div className="team-stat-item">
                                                <span className="number">{team.task_count}</span>
                                                <span className="label">Tasks</span>
                                            </div>
                                            <div className="team-stat-item">
                                                <span className="number">{team.member_count}</span>
                                                <span className="label">Members</span>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </section>
                    )}

                    {/* Recent Tasks */}
                    <section className="dashboard-section animate-fadeIn">
                        <div className="section-header">
                            <h2>Recent Tasks</h2>
                            <Link to="/board" className="btn btn-primary">
                                + New Task
                            </Link>
                        </div>

                        <div className="tasks-list">
                            {issues.length > 0 ? (
                                issues.slice(0, 5).map(issue => (
                                    <div key={issue.id} className="task-item">
                                        <div className="task-info">
                                            <h4>{issue.title}</h4>
                                            <div className="task-meta">
                                                <span
                                                    className="task-status"
                                                    style={{ backgroundColor: issue.status?.color || '#6B7280' }}
                                                >
                                                    {issue.status?.name || 'No Status'}
                                                </span>
                                                <span className={`badge ${getPriorityClass(issue.priority)}`}>
                                                    {issue.priority}
                                                </span>
                                            </div>
                                        </div>
                                        <Link to="/board" className="btn btn-secondary">View</Link>
                                    </div>
                                ))
                            ) : (
                                <div className="empty-state">
                                    <span className="empty-icon">📭</span>
                                    <h3>No tasks yet</h3>
                                    <p>Create your first task to get started</p>
                                    <Link to="/board" className="btn btn-primary" style={{ marginTop: '1rem' }}>
                                        Go to Board
                                    </Link>
                                </div>
                            )}
                        </div>
                    </section>
                </div>

                {/* Quick Info */}
                <section className="dashboard-section">
                    <h2>Account Info</h2>
                    <div className="info-grid">
                        <div className="info-card">
                            <span className="info-label">Email</span>
                            <span className="info-value">{user?.email}</span>
                        </div>
                        <div className="info-card">
                            <span className="info-label">Organization</span>
                            <span className="info-value">{orgName || 'Loading...'}</span>
                        </div>
                        <div className="info-card">
                            <span className="info-label">Teams</span>
                            <span className="info-value">{teams.length}</span>
                        </div>
                    </div>
                    {permissions && permissions.teams.length > 0 && (
                        <div className="roles-section">
                            <h4>My Roles</h4>
                            <div className="roles-list">
                                {permissions.teams.map(t => (
                                    <div key={t.team_id} className="role-item">
                                        <span className="role-team">{t.team_name}</span>
                                        <span className={`role-badge ${t.role}`}>{t.role}</span>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}
                </section>
            </main>
        </div>
    );
};

export default Dashboard;
