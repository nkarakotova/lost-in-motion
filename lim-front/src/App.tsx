import React, { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import SignupPage from './pages/SignupPage';
import LoginPage from './pages/LoginPage';
import SchedulePage from './pages/SchedulePage';
import ChangePasswordPage from './pages/ChangePasswordPage';
import CreateTraining from './pages/CreateTrainingPage';
import NotFound from './pages/NotFound';
import { Configuration } from './api/client/configuration';
import { ThemeProvider } from '@mui/material/styles';
import theme from './styles/theme';

const App: React.FC = () => {
  useEffect(() => {
    const config = new Configuration();
    config.addAuthorizationHeader();
  }, []);

  return (
    <ThemeProvider theme={theme}>
      <Router>
        <Routes>
          <Route path="/signup" element={<SignupPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/schedule" element={<SchedulePage />} />
          <Route path="/change-password" element={<ChangePasswordPage />} />
          <Route path="/create-training" element={<CreateTraining />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
};

export default App;
