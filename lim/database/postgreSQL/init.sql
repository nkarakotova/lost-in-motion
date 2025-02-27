drop table if exists coaches cascade;
create table coaches
(
	coach_id serial primary key,
    name text
);
alter table coaches add constraint unique_name_coach unique (name);

drop table if exists halls cascade;
create table halls
(
	hall_id serial primary key,
    number int
);
alter table halls add constraint unique_number_hall unique (number);

drop table if exists clients cascade;
create table clients
(
	client_id serial primary key,
    name text,
	telephone text,
	mail text,
    password text
);
alter table clients add constraint unique_telephone_client unique (telephone);

drop table if exists trainings cascade;
create table trainings 
(
	training_id serial  primary key,
	coach_id int references coaches (coach_id),
	hall_id int references halls (hall_id),
    name text,
	date_time timestamp,
	places_num int
);

drop table if exists clients_trainings cascade;
create table clients_trainings
(
	client_id int references clients(client_id) on delete cascade, 
	training_id int references trainings(training_id) on delete cascade
);
