import React, { useState } from 'react';
import { AppBar, Toolbar, Button, Typography, Box, Dialog, DialogActions, DialogContent, DialogTitle, TextField, CircularProgress } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { CoachService } from '../api/services/coachesService';
import { HallService } from '../api/services/hallService';
import { CoachInfo } from '../api/client/api';
import { HallInfo } from '../api/client/api';
import Notification from '../components/Notification'; 

const getUserRole = () => {
  const token = localStorage.getItem('jwt_token');
  if (!token) return null;
  const decoded: any = JSON.parse(atob(token.split('.')[1]));
  return decoded.role;
};

const Header: React.FC = () => {
  const navigate = useNavigate();
  const role = getUserRole();

  const [openCoachDialog, setOpenCoachDialog] = useState(false);
  const [openHallDialog, setOpenHallDialog] = useState(false);
  const [coachName, setCoachName] = useState<string>('');
  const [hallNumber, setHallNumber] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false); 
  const [openSnackbar, setOpenSnackbar] = useState<boolean>(false);
  const [snackbarSeverity, setSnackbarSeverity] = useState<'success' | 'error'>('success'); 
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  
  
  const handleOpenCoachDialog = () => setOpenCoachDialog(true);
  const handleOpenHallDialog = () => setOpenHallDialog(true);

  
  const handleCloseCoachDialog = () => {
    setOpenCoachDialog(false);
    setCoachName('');
  };

  const handleCloseHallDialog = () => {
    setOpenHallDialog(false);
    setHallNumber('');
  };

  const handleCreateCoach = async () => {
    if (!coachName) {
      setError('Необходимо ввести имя тренера!');
      setSnackbarSeverity('error');
      setOpenSnackbar(true);
      return;
    }

    const coachInfo: CoachInfo = {
      name: coachName,
    };

    setLoading(true);
    const coachService = new CoachService();

    try {
      await coachService.createCoach(coachInfo, setError);
      setSuccessMessage('Тренер успешно добавлен!');
      setError(null);
      setOpenSnackbar(true); 
    } catch (err: any) {
      setSuccessMessage(null);
      setError(err);
      setOpenSnackbar(true);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateHall = async () => {
    if (!hallNumber) {
      setError('Пожалуйста, укажите номер зала');
      setOpenSnackbar(true);
      return;
    }

    const hallNumberValue = Number(hallNumber);
    if (isNaN(hallNumberValue)) {
      setError('Номер зала должен быть числом');
      setOpenSnackbar(true);
      return;
    }

    const hallInfo: HallInfo = {
      number: hallNumberValue,
    };

    setLoading(true);
    const hallService = new HallService();

    try {
      await hallService.createHall(hallInfo, setError);
      setSuccessMessage('Зал успешно добавлен!');
      setError(null);
      setOpenSnackbar(true);
    } catch (err: any) {
      setSuccessMessage(null);
      setError(err);
      setOpenSnackbar(true);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('jwt_token');
    navigate('/login');
  };

  const renderButtons = () => {
    if (!role) {
      return (
        <>
          <Button color="inherit" onClick={() => navigate('/login')}>
            <Typography variant="body1">Вход</Typography>
          </Button>
          <Button color="inherit" onClick={() => navigate('/signup')}>
            <Typography variant="body1">Регистрация</Typography>
          </Button>
        </>
      );
    } else if (role === 'client') {
      return (
        <>
          <Button color="inherit" onClick={() => navigate('/change-password')}>
            <Typography variant="body1">Изменить пароль</Typography>
          </Button>
          <Button color="inherit" onClick={handleLogout}>
            <Typography variant="body1">Выход</Typography>
          </Button>
        </>
      );
    } else if (role === 'admin') {
      return (
        <>
          <Button color="inherit" onClick={() => navigate('/create-training')}>
            <Typography variant="body1">Добавить тренировку</Typography>
          </Button>
          <Button color="inherit" onClick={handleOpenHallDialog}>
            <Typography variant="body1">Добавить зал</Typography>
          </Button>
          <Button color="inherit" onClick={handleOpenCoachDialog}>
            <Typography variant="body1">Добавить тренера</Typography>
          </Button>
          <Button color="inherit" onClick={handleLogout}>
            <Typography variant="body1">Выход</Typography>
          </Button>
        </>
      );
    }
    return null;
  };

  return (
    <AppBar position="static" color="primary">
      <Toolbar>
        <Button color="inherit" onClick={() => navigate('/schedule')}>
          <Typography variant="body1">Домой</Typography>
        </Button>
        <Box flexGrow={1} display="flex" justifyContent="center">
          <Typography variant="h1" component="div" sx={{ position: 'absolute', left: '50%', top: '20%', transform: 'translate(-50%, -10%)', textAlign: 'center' }}>
            Расписание
          </Typography>
        </Box>
        {renderButtons()}
      </Toolbar>

      {/* Диалоги для добавления тренера и зала */}
      <Dialog open={openCoachDialog} onClose={handleCloseCoachDialog}>
        <DialogTitle>Добавить тренера</DialogTitle>
        <DialogContent>
          <TextField
            label="Имя тренера"
            fullWidth
            value={coachName}
            onChange={(e) => setCoachName(e.target.value)}
            margin="normal"
            variant="outlined"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseCoachDialog} color="primary">Отмена</Button>
          <Button onClick={handleCreateCoach} color="primary" disabled={loading}>
            {loading ? <CircularProgress size={24} /> : 'Добавить'}
          </Button>
        </DialogActions>
      </Dialog>

      <Dialog open={openHallDialog} onClose={handleCloseHallDialog}>
        <DialogTitle>Добавить зал</DialogTitle>
        <DialogContent>
          <TextField
            label="Номер зала"
            fullWidth
            value={hallNumber}
            onChange={(e) => setHallNumber(e.target.value)}
            margin="normal"
            variant="outlined"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseHallDialog} color="primary">Отмена</Button>
          <Button onClick={handleCreateHall} color="primary" disabled={loading}>
            {loading ? <CircularProgress size={24} /> : 'Добавить'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Всплывающее окно для ошибок и успеха */}
      <Notification
        open={openSnackbar}
        message={error || successMessage}
        severity={error ? 'error' : 'success'}
        onClose={() => setOpenSnackbar(false)}
      />
    </AppBar>
  );
};

export default Header;
