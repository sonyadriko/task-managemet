import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { Link } from 'react-router-dom';
import apiClient from '../api/client';
import Sidebar from '../components/Sidebar';
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

const Dashboard: React.FC = () => {
    const { user } = useAuth();
    const [teams, setTeams] = useState<Team[]>([]);
    const [issues, setIssues] = useState<Issue[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const teamsRes = await apiClient.get('/teams');
                setTeams(teamsRes.data || []);

                if (teamsRes.data && teamsRes.data.length > 0) {
                    const issuesRes = await apiClient.get(`/issues?team_id=${teamsRes.data[0].id}`);
                    setIssues(issuesRes.data || []);
                }
            } catch (error) {
                console.error('Failed to fetch data:', error);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, []);

    const getPriorityClass = (priority: string) => {
        switch (priority) {
            case 'URGENT': return 'badge-danger';
            case 'HIGH': return 'badge-warning';
            case 'NORMAL': return 'badge-primary';
            default: return 'badge-secondary';
        }
    };

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
                            <h3>{issues.length}</h3>
                            <p>Total Tasks</p>
                        </div>
                    </div>
                    <div className="stat-card">
                        <div className="stat-icon blue">🔄</div>
                        <div className="stat-content">
                            <h3>{issues.filter(i => i.status?.name === 'IN_PROGRESS').length}</h3>
                            <p>In Progress</p>
                        </div>
                    </div>
                    <div className="stat-card">
                        <div className="stat-icon green">✅</div>
                        <div className="stat-content">
                            <h3>{issues.filter(i => i.status?.name === 'DONE').length}</h3>
                            <p>Completed</p>
                        </div>
                    </div>
                    <div className="stat-card">
                        <div className="stat-icon orange">👥</div>
                        <div className="stat-content">
                            <h3>{teams.length}</h3>
                            <p>Teams</p>
                        </div>
                    </div>
                </div>

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

                {/* Quick Info */}
                <section className="dashboard-section">
                    <h2>Account Info</h2>
                    <div className="info-grid">
                        <div className="info-card">
                            <span className="info-label">Email</span>
                            <span className="info-value">{user?.email}</span>
                        </div>
                        <div className="info-card">
                            <span className="info-label">Organization ID</span>
                            <span className="info-value">{user?.organization_id}</span>
                        </div>
                        <div className="info-card">
                            <span className="info-label">Timezone</span>
                            <span className="info-value">{user?.timezone}</span>
                        </div>
                    </div>
                </section>
            </main>
        </div>
    );
};

export default Dashboard;
