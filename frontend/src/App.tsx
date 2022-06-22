import React, { useEffect } from 'react';
import { Navigate, Route, Routes, useLocation, useSearchParams } from 'react-router-dom';
import { Home } from './pages/home/Home';
import { LoginPage } from './pages/login/LoginPage';

export function App() {


  const [searchParams, setSearchParams] = useSearchParams();

  useEffect(() => {
    const accessToken = searchParams.get('accessToken');
    const refreshToken = searchParams.get('refreshToken');
    if (accessToken && refreshToken) {
      localStorage.setItem('accessToken', accessToken);
      localStorage.setItem('refreshToken', refreshToken);
      searchParams.delete('accessToken');
      searchParams.delete('refreshToken');
      setSearchParams(searchParams);
    }
  });

  return (
    <AuthProvider>
      <Routes>
        <Route path="/" element={<RequireAuth><Home/></RequireAuth>}></Route>
        <Route path="/login" element={<LoginPage/>}/>
      </Routes>
    </AuthProvider>
  );
}

interface AuthContextType {
  accessToken: string | null;
}

const AuthContext = React.createContext<AuthContextType>(null!);

function AuthProvider({ children }: { children: React.ReactNode }) {

  const value = {
    accessToken: localStorage.getItem('accessToken'),
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

function useAuth() {
  return React.useContext(AuthContext);
}

function RequireAuth({ children }: { children: JSX.Element }) {
  const auth = useAuth();
  const location = useLocation();

  if (!auth.accessToken) {
    return <Navigate to="/login" state={{ from: location }} replace/>;
  }

  return children;
}
