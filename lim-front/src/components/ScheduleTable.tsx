import React, { useState, useEffect } from 'react';
import { Training } from '../api/client/api';
import { Table, TableBody, TableCell, TableContainer, Typography, TableHead, TableRow, CircularProgress, Paper, Button, Dialog, DialogActions, DialogContent, DialogTitle } from '@mui/material';
import { useTheme } from '@mui/material/styles';
import { useNavigate } from 'react-router-dom';
import ActionButton from './ActionButton';
import { ClientService } from '../api/services/clientService';
import TrainingAssignmentService from '../api/services/assignmentService';
import { TrainingService } from '../api/services/trainingsService';
import Notification from './Notification';

interface ScheduleTableProps {
  fetchData: () => Promise<Training[]>;
}

const ScheduleTable: React.FC<ScheduleTableProps> = ({ fetchData }) => {
  const [schedule, setSchedule] = useState<Training[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [snackbarOpen, setSnackbarOpen] = useState<boolean>(false);
  const theme = useTheme();
  const navigate = useNavigate();
  const [clientTrainings, setClientTrainings] = useState<Training[]>([]);

  const [openDialog, setOpenDialog] = useState(false);
  const [currentAction, setCurrentAction] = useState<'register' | 'cancel' | 'delete' | null>(null);
  const [selectedTraining, setSelectedTraining] = useState<Training | null>(null);

  const getUserRole = () => {
    const token = localStorage.getItem('jwt_token');
    if (!token) return null;

    const decoded: any = JSON.parse(atob(token.split('.')[1]));
    return decoded.role;
  };

  const userRole = getUserRole();

  useEffect(() => {
    const fetchSchedule = async () => {
      try {
        const data = await fetchData();
        const sortedData = data.sort((a, b) => {
          const dateA = new Date(a.date_time).getTime();
          const dateB = new Date(b.date_time).getTime();
          return dateA - dateB;
        });
        setSchedule(sortedData);

        if (userRole === 'client') {
          const clientService = new ClientService();
          const response = await clientService.getTrainingsByClient(setError);
          setClientTrainings(response);
        }
      } catch (err: any) {
        setError(err);
      } finally {
        setLoading(false);
      }
    };

    fetchSchedule();
  }, [fetchData, userRole]);

  if (loading) return <CircularProgress />;

  const renderActionButton = (item: Training) => {
    if (!userRole) {
      return <ActionButton label="Записаться" onClick={() => navigate('/signup')} />;
    }

    if (userRole === 'admin') {
      return <ActionButton label="Удалить" onClick={() => handleOpenDialog('delete', item)} />;
    }

    if (userRole === 'client') {
      const isRegistered = clientTrainings.some(training => training.training_id === item.training_id);

      if (isRegistered) {
        return <ActionButton label="Удалить запись" onClick={() => handleOpenDialog('cancel', item)} />;
      } else {
        return <ActionButton label="Записаться" onClick={() => handleOpenDialog('register', item)} />;
      }
    }

    return null;
  };

  const handleOpenDialog = (action: 'register' | 'cancel' | 'delete', item: Training) => {
    setCurrentAction(action);
    setSelectedTraining(item);
    setOpenDialog(true);
  };

  const handleCloseDialog = () => {
    setOpenDialog(false);
    setCurrentAction(null);
    setSelectedTraining(null);
  };

  const handleConfirmAction = async () => {
    if (!selectedTraining) return;

    const trainingAssignmentService = new TrainingAssignmentService();
    const trainingService = new TrainingService();

    try {
      setError(null);
      setSuccessMessage(null);

      if (currentAction === 'register') {
        await trainingAssignmentService.createAssignment({ training_id: selectedTraining.training_id }, setError);

        setSuccessMessage('Вы успешно записались на тренировку!');
        setClientTrainings([...clientTrainings, selectedTraining]);
      } else if (currentAction === 'cancel') {
        await trainingAssignmentService.deleteAssignment({ training_id: selectedTraining.training_id }, setError);
        setSuccessMessage('Вы успешно удалили запись на тренировку.');
        setClientTrainings(clientTrainings.filter(training => training.training_id !== selectedTraining.training_id));
      } else if (currentAction === 'delete') {
        await trainingService.deleteTraining({ training_id: selectedTraining.training_id }, setError);
        setSuccessMessage('Вы успешно удалили тренировку.');
        setSchedule(schedule.filter(training => training.training_id !== selectedTraining.training_id));
      }

      setError(null);
      setSnackbarOpen(true);
      handleCloseDialog();
    } catch (err: any) {
      setSuccessMessage(null);
      setError(err);

      setSnackbarOpen(true);
      handleCloseDialog();
    }
  };

  return (
    <div style={{ padding: '20px' }}>
      <TableContainer component={Paper} sx={{ backgroundColor: theme.palette.background.default }}>
        <Table aria-label="schedule table">
          <TableHead>
            <TableRow>
              {['Название тренировки', 'Дата', 'Время', 'Количество мест', 'Зал', 'Тренер', 'Действие'].map((header) => (
                <TableCell
                  key={header}
                  sx={{
                    ...theme.typography.h2,
                    color: theme.palette.text.primary,
                    backgroundColor: theme.palette.background.paper,
                    textAlign: 'center',
                  }}
                >
                  <strong>{header}</strong>
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {schedule.map((item) => {
              const trainingDate = new Date(item.date_time);
              
              const formattedDate = trainingDate.toLocaleDateString();
              const formattedTime = trainingDate.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });

              
              
              trainingDate.setHours(trainingDate.getHours() - 3);

              return (
                <TableRow key={item.training_id}>
                  <TableCell sx={{ ...theme.typography.body1, color: theme.palette.text.primary, textAlign: 'center' }}>
                    {item.name}
                  </TableCell>
                  <TableCell sx={{ ...theme.typography.body1, color: theme.palette.text.primary, textAlign: 'center' }}>
                    {formattedDate}
                  </TableCell>
                  <TableCell sx={{ ...theme.typography.body1, color: theme.palette.text.primary, textAlign: 'center' }}>
                    {trainingDate.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                  </TableCell>
                  <TableCell sx={{ ...theme.typography.body1, color: theme.palette.text.primary, textAlign: 'center' }}>
                    {item.places_num}
                  </TableCell>
                  <TableCell sx={{ ...theme.typography.body1, color: theme.palette.text.primary, textAlign: 'center' }}>
                    {item.hall_number}
                  </TableCell>
                  <TableCell sx={{ ...theme.typography.body1, color: theme.palette.text.primary, textAlign: 'center' }}>
                    {item.coach_name}
                  </TableCell>
                  <TableCell sx={{ ...theme.typography.body1, color: theme.palette.text.primary, textAlign: 'center' }}>
                    {renderActionButton(item)}
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </TableContainer>

      <Notification
        open={snackbarOpen}
        message={error || successMessage}
        severity={error ? 'error' : 'success'}
        onClose={() => setSnackbarOpen(false)}
      />

      <Dialog open={openDialog} onClose={handleCloseDialog}>
        <DialogTitle>Подтверждение</DialogTitle>
        <DialogContent>
          <Typography>
            {currentAction === 'register'
              ? 'Вы уверены, что хотите записаться на эту тренировку?'
              : currentAction === 'cancel'
              ? 'Вы уверены, что хотите отменить запись на тренировку?'
              : 'Вы уверены, что хотите удалить эту тренировку?'}
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDialog} color="primary">
            Отмена
          </Button>
          <Button onClick={handleConfirmAction} color="primary">
            Подтвердить
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default ScheduleTable;
