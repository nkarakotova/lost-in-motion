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
        (1, 1, 'Clasic', '2024-12-17 11:00:00', 3),
        (1, 2, 'Tango', '2025-03-15 11:00:00', 8),
        (2, 1, 'Hiphop', '2025-07-20 17:00:00', 12),
        (3, 3, 'Tango', '2025-11-10 11:00:00', 7),
        (1, 3, 'Clasic',        '2026-02-14 11:00:00', 5),
        (1, 1, 'Tango',         '2026-05-12 11:00:00', 8),
        (2, 2, 'Hiphop',        '2026-05-13 17:00:00', 10),
        (3, 3, 'Contemporary',  '2026-05-14 11:00:00', 6),
        (1, 2, 'Clasic',        '2026-05-15 18:00:00', 5),
        (2, 1, 'Hiphop',        '2026-05-16 17:00:00', 12),
        (3, 2, 'Tango',         '2026-05-17 11:00:00', 7),
        (1, 1, 'Tango',         '2026-06-15 11:00:00', 8),
        (2, 2, 'Hiphop',        '2026-08-20 17:00:00', 10),
        (3, 1, 'Contemporary',  '2026-09-01 11:00:00', 6);

insert into clients_trainings (client_id, training_id)
values  (1, 2),
        (2, 3),
        (3, 1),
        (1, 7),
        (2, 7),
        (3, 8),
        (1, 9),
        (2, 10),
        (1, 14),
        (2, 14),
        (3, 15),
        (1, 16);
