import React from 'react';
import { useAuth } from '../contexts/AuthContext';
import Sidebar from '../components/Sidebar';
import './Settings.css';

const Settings: React.FC = () => {
    const { user } = useAuth();

    return (
        <div className="page-layout">
            <Sidebar activeItem="settings" />
            <main className="main-content">
                <header className="page-header">
                    <div className="header-left">
                        <h1>Settings</h1>
                        <p>Manage your account and preferences</p>
                    </div>
                </header>

                <div className="settings-grid">
                    {/* Profile Section */}
                    <section className="settings-section">
                        <h2>👤 Profile</h2>
                        <div className="settings-card">
                            <div className="profile-header">
                                <div className="profile-avatar">
                                    {user?.full_name?.charAt(0) || 'U'}
                                </div>
                                <div className="profile-info">
                                    <h3>{user?.full_name}</h3>
                                    <p>{user?.email}</p>
                                </div>
                            </div>
                            <div className="form-group">
                                <label className="form-label">Full Name</label>
                                <input type="text" className="form-input" value={user?.full_name || ''} readOnly />
                            </div>
                            <div className="form-group">
                                <label className="form-label">Email</label>
                                <input type="email" className="form-input" value={user?.email || ''} readOnly />
                            </div>
                            <div className="form-group">
                                <label className="form-label">Timezone</label>
                                <input type="text" className="form-input" value={user?.timezone || ''} readOnly />
                            </div>
                            <button className="btn btn-secondary" disabled>
                                Update Profile (Coming Soon)
                            </button>
                        </div>
                    </section>

                    {/* Organization Section */}
                    <section className="settings-section">
                        <h2>🏢 Organization</h2>
                        <div className="settings-card">
                            <div className="info-row">
                                <span className="info-label">Organization ID</span>
                                <span className="info-value">{user?.organization_id}</span>
                            </div>
                            <div className="info-row">
                                <span className="info-label">Role</span>
                                <span className="info-value badge badge-primary">Member</span>
                            </div>
                        </div>
                    </section>

                    {/* Security Section */}
                    <section className="settings-section">
                        <h2>🔒 Security</h2>
                        <div className="settings-card">
                            <div className="form-group">
                                <label className="form-label">Current Password</label>
                                <input type="password" className="form-input" placeholder="••••••••" disabled />
                            </div>
                            <div className="form-group">
                                <label className="form-label">New Password</label>
                                <input type="password" className="form-input" placeholder="••••••••" disabled />
                            </div>
                            <div className="form-group">
                                <label className="form-label">Confirm New Password</label>
                                <input type="password" className="form-input" placeholder="••••••••" disabled />
                            </div>
                            <button className="btn btn-secondary" disabled>
                                Change Password (Coming Soon)
                            </button>
                        </div>
                    </section>

                    {/* Preferences Section */}
                    <section className="settings-section">
                        <h2>🎨 Preferences</h2>
                        <div className="settings-card">
                            <div className="preference-row">
                                <div>
                                    <h4>Dark Mode</h4>
                                    <p>Use dark theme for the interface</p>
                                </div>
                                <label className="toggle">
                                    <input type="checkbox" checked readOnly />
                                    <span className="toggle-slider"></span>
                                </label>
                            </div>
                            <div className="preference-row">
                                <div>
                                    <h4>Email Notifications</h4>
                                    <p>Receive email updates for task changes</p>
                                </div>
                                <label className="toggle">
                                    <input type="checkbox" readOnly />
                                    <span className="toggle-slider"></span>
                                </label>
                            </div>
                        </div>
                    </section>
                </div>
            </main>
        </div>
    );
};

export default Settings;
