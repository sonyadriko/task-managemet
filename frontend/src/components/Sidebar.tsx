import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import apiClient from '../api/client';
import './Sidebar.css';

interface Team {
    id: number;
    name: string;
}

interface SidebarProps {
    activeItem: string;
}

const Sidebar: React.FC<SidebarProps> = ({ activeItem }) => {
    const { logout } = useAuth();
    const navigate = useNavigate();
    const [teams, setTeams] = useState<Team[]>([]);

    useEffect(() => {
        const fetchTeams = async () => {
            try {
                const res = await apiClient.get('/teams');
                setTeams(res.data || []);
            } catch (error) {
                console.error('Failed to fetch teams:', error);
            }
        };
        fetchTeams();
    }, []);

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    const navItems = [
        { id: 'dashboard', icon: '🏠', label: 'Dashboard', path: '/dashboard' },
        { id: 'board', icon: '📊', label: 'Board', path: '/board' },
        { id: 'calendar', icon: '📅', label: 'Calendar', path: '/calendar' },
        { id: 'teams', icon: '👥', label: 'Teams', path: '/teams' },
        { id: 'settings', icon: '⚙️', label: 'Settings', path: '/settings' },
    ];

    return (
        <aside className="sidebar">
            <div className="sidebar-header">
                <div className="sidebar-logo">
                    <span>🌊</span>
                </div>
                <h2>Flow</h2>
            </div>

            <nav className="sidebar-nav">
                {navItems.map(item => (
                    <Link
                        key={item.id}
                        to={item.path}
                        className={`nav-item ${activeItem === item.id ? 'active' : ''}`}
                    >
                        <span className="nav-icon">{item.icon}</span>
                        {item.label}
                    </Link>
                ))}
            </nav>

            <div className="sidebar-teams">
                <h3>Your Teams</h3>
                {teams.map(team => (
                    <div key={team.id} className="team-item">
                        <span className="team-avatar">{team.name.charAt(0)}</span>
                        <span className="team-name">{team.name}</span>
                    </div>
                ))}
                {teams.length === 0 && (
                    <p className="no-teams">No teams yet</p>
                )}
            </div>

            <div className="sidebar-footer">
                <button onClick={handleLogout} className="logout-btn">
                    <span>🚪</span> Logout
                </button>
            </div>
        </aside>
    );
};

export default Sidebar;
