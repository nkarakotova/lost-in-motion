insert into coaches (name)
values  ('Kate'),
        ('Denis'),
        ('Alis');

insert into halls (number)
values  (1),
        (3),
        (5);

insert into clients (name, telephone, mail, password)
values  ('Diana', '1111111111', 'd@mail.ru', '111'),
        ('Maria', '2222222222', 'm@mail.ru', '222'),
        ('Anton', '3333333333', 'a@mail.ru', '333');

insert into trainings (coach_id, hall_id, name, date_time, places_num)
values  (2, 1, 'Hiphop', '2024-12-05 17:00:00', 11),
        (1, 3, 'Clasic', '2024-12-05 17:00:00', 1),
        (1, 2, 'Tango', '2024-12-05 11:00:00', 7),
        (3, 2, 'Tango', '2024-12-07 11:00:00', 7),
        (1, 1, 'Clasic', '2024-12-07 11:00:00', 3),
        (1, 1, 'Clasic', '2024-12-17 11:00:00', 3);

insert into clients_trainings (client_id, training_id)
values  (1, 2),
        (2, 3),
        (3, 1);
