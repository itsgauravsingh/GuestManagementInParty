create database if not exists partydb;
use partydb;

create table if not exists Guest
(
	id int auto_increment
		primary key,
	guestname varchar(20) not null,
	accompanying_guests int not null,
	tableId int not null,
	constraint Guest_name_uindex
		unique (guestname)
);

create table if not exists GuestLog
(
	id int auto_increment
		primary key,
	guestid varchar(20) not null,
	ispresent int default 1 not null,
	time_arrived varchar(30) not null,
	time_departed varchar(30) null,
	accompanying_guest int default 0 null,
	constraint GuestLog_pk
		unique (guestid, ispresent)
)
comment 'this will contain the data about who came, when came, with how many, when departed';

create table if not exists TableInfo
(
	id int not null,
	state varchar(16) default 'available' null,
	capacity int default 10 not null,
	venueid int default 1 not null,
	primary key (id, venueid)
)
comment 'used to store table entities';

create index TableInfo_id_index
	on TableInfo (id);

create table if not exists Venue
(
	id int not null
		primary key,
	venuename varchar(11) not null
)
comment 'this is used to keep venue information';

insert into Guest (guestname, accompanying_guests, tableId)
values ('Jemima Clark', 4,2), ('Vera Yep',5,4),('Gaurav Singh',2,3),('Jeff Bezos',10,1),('Moubin Faizullah-Khan',20,5);

insert into TableInfo (id, capacity, venueid)
values (1,15,1), (2,5,1),(3,2,1),(4,8,1),(5,25,1),(6,5,1),(7,5,1);