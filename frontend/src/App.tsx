import { Routes, Route, Navigate, Link, useLocation, useNavigate } from 'react-router-dom'
import { isAuthenticated, clearToken } from './api/client'
import InvitePage from './pages/InvitePage'
import AdminLogin from './pages/AdminLogin'
import AdminDashboard from './pages/AdminDashboard'
import AdminQRScanner from './pages/AdminQRScanner'

function AdminLayout() {
  const navigate = useNavigate()
  const location = useLocation()

  function handleLogout() {
    clearToken()
    navigate('/admin/login')
  }

  return (
    <div className="admin-layout">
      <aside className="admin-sidebar">
        <div className="brand">
          <div className="brand-icon">◆</div>
          <div className="brand-text">The <span>Wedding</span></div>
        </div>
        <nav>
          <Link
            to="/admin/dashboard"
            className={location.pathname === '/admin/dashboard' ? 'active' : ''}
          >
            Dashboard
          </Link>
          <Link
            to="/admin/qr_scanner"
            className={location.pathname === '/admin/qr_scanner' ? 'active' : ''}
          >
            QR Scanner
          </Link>
        </nav>
        <div className="sidebar-footer">
          <button className="logout-btn" onClick={handleLogout}>
            Log out
          </button>
        </div>
      </aside>
      {location.pathname === '/admin/qr_scanner' ? <AdminQRScanner /> : <AdminDashboard />}
    </div>
  )
}

function RequireAuth({ children }: { children: React.ReactNode }) {
  if (!isAuthenticated()) {
    return <Navigate to="/admin/login" replace />
  }
  return <>{children}</>
}

export default function App() {
  return (
    <Routes>
      <Route path="/invite/:external_id" element={<InvitePage />} />
      <Route path="/admin/login" element={<AdminLogin />} />
      <Route
        path="/admin/dashboard"
        element={
          <RequireAuth>
            <AdminLayout />
          </RequireAuth>
        }
      />
      <Route
        path="/admin/qr_scanner"
        element={
          <RequireAuth>
            <AdminLayout />
          </RequireAuth>
        }
      />
      <Route path="*" element={<Navigate to="/admin/login" replace />} />
    </Routes>
  )
}
