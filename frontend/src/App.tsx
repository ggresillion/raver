import React from 'react';
import { Route, Routes, useSearchParams } from 'react-router-dom';
import { Home } from './pages/home/Home';

export function App() {

  const [searchParams, setSearchParams] = useSearchParams();
  const accessToken = searchParams.get('accessToken');
  const refreshToken = searchParams.get('refreshToken');
  if (accessToken && refreshToken) {
    localStorage.setItem('accessToken', accessToken);
    localStorage.setItem('refreshToken', refreshToken);
    searchParams.delete('accessToken');
    searchParams.delete('refreshToken');
    setSearchParams(searchParams);
  }

  return (
    <Routes>
      <Route path='/' element={<Home />}></Route>
    </Routes>
  );
}
