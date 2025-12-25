import React, { useState, useEffect } from 'react';
import apiClient from '../api/client';
import Sidebar from '../components/Sidebar';
import './Calendar.css';

interface CalendarEvent {
    issue_id: number;
    issue_title: string;
    priority: string;
    status_name: string;
    status_color: string;
    user_name: string;
    team_name: string;
    start_date: string;
    end_date: string;
}

interface Team {
    id: number;
    name: string;
}

const Calendar: React.FC = () => {
    const [events, setEvents] = useState<CalendarEvent[]>([]);
    const [teams, setTeams] = useState<Team[]>([]);
    const [selectedTeam, setSelectedTeam] = useState<number | null>(null);
    const [currentDate, setCurrentDate] = useState(new Date());
    const [loading, setLoading] = useState(true);

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

        const fetchEvents = async () => {
            setLoading(true);
            const startDate = new Date(currentDate.getFullYear(), currentDate.getMonth(), 1);
            const endDate = new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 0);

            try {
                const res = await apiClient.get('/calendar', {
                    params: {
                        team_id: selectedTeam,
                        start_date: startDate.toISOString().split('T')[0],
                        end_date: endDate.toISOString().split('T')[0]
                    }
                });
                setEvents(res.data || []);
            } catch (error) {
                console.error('Failed to fetch events:', error);
            } finally {
                setLoading(false);
            }
        };
        fetchEvents();
    }, [selectedTeam, currentDate]);

    const getDaysInMonth = () => {
        const year = currentDate.getFullYear();
        const month = currentDate.getMonth();
        const firstDay = new Date(year, month, 1);
        const lastDay = new Date(year, month + 1, 0);
        const daysInMonth = lastDay.getDate();
        const startingDay = firstDay.getDay();

        const days = [];

        // Previous month padding
        for (let i = 0; i < startingDay; i++) {
            days.push({ day: null, date: null });
        }

        // Current month days
        for (let i = 1; i <= daysInMonth; i++) {
            days.push({
                day: i,
                date: new Date(year, month, i)
            });
        }

        return days;
    };

    const getEventsForDay = (date: Date | null) => {
        if (!date) return [];
        const dateStr = date.toISOString().split('T')[0];
        return events.filter(event => {
            const start = event.start_date.split('T')[0];
            const end = event.end_date.split('T')[0];
            return dateStr >= start && dateStr <= end;
        });
    };

    const prevMonth = () => {
        setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() - 1, 1));
    };

    const nextMonth = () => {
        setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 1));
    };

    const formatMonth = () => {
        return currentDate.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
    };

    const isToday = (date: Date | null) => {
        if (!date) return false;
        const today = new Date();
        return date.toDateString() === today.toDateString();
    };

    const days = getDaysInMonth();
    const weekDays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

    return (
        <div className="page-layout">
            <Sidebar activeItem="calendar" />
            <main className="main-content">
                <header className="page-header">
                    <div className="header-left">
                        <h1>Calendar</h1>
                        <p>View task timeline and assignments</p>
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
                    </div>
                </header>

                <div className="calendar-nav">
                    <button className="btn btn-secondary" onClick={prevMonth}>← Prev</button>
                    <h2>{formatMonth()}</h2>
                    <button className="btn btn-secondary" onClick={nextMonth}>Next →</button>
                </div>

                {loading ? (
                    <div className="loading-screen">
                        <div className="spinner-lg"></div>
                    </div>
                ) : (
                    <div className="calendar-grid">
                        {/* Week day headers */}
                        {weekDays.map(day => (
                            <div key={day} className="calendar-header">{day}</div>
                        ))}

                        {/* Calendar days */}
                        {days.map((item, index) => (
                            <div
                                key={index}
                                className={`calendar-day ${!item.day ? 'empty' : ''} ${isToday(item.date) ? 'today' : ''}`}
                            >
                                {item.day && (
                                    <>
                                        <span className="day-number">{item.day}</span>
                                        <div className="day-events">
                                            {getEventsForDay(item.date).slice(0, 3).map((event, i) => (
                                                <div
                                                    key={i}
                                                    className="event-item"
                                                    style={{ backgroundColor: event.status_color }}
                                                    title={`${event.issue_title} - ${event.user_name}`}
                                                >
                                                    {event.issue_title.substring(0, 15)}...
                                                </div>
                                            ))}
                                            {getEventsForDay(item.date).length > 3 && (
                                                <div className="more-events">
                                                    +{getEventsForDay(item.date).length - 3} more
                                                </div>
                                            )}
                                        </div>
                                    </>
                                )}
                            </div>
                        ))}
                    </div>
                )}
            </main>
        </div>
    );
};

export default Calendar;
