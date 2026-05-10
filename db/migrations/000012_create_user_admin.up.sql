insert into users(firstname, lastname, birthday, phone_number, password, role) select 'admin','admin', '2026-05-09T14:32:47.123456+05:00','77701320091', '$2a$10$wly60YDSXJYe3EwWV3gxYeEHMeuvA9ZM5688Cs2uGR9B1te3diN1q', 'admin' where not exists (select 1 from users where phone_number = '77701320091');






