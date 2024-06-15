--photo_types
INSERT INTO photo_types(id, name) VALUES (1, 'С использованием БПЛА');
INSERT INTO photo_types(id, name) VALUES (2, 'Проведение профессиональной кино-, фото- и видеосъемки со стационарным оборудованием');
INSERT INTO photo_types(id, name) VALUES (3, 'Без съемки');


--visit_reasons
INSERT INTO visit_reasons(id, name) VALUES (1, 'Однодневная поездка/экскурсия');
INSERT INTO visit_reasons(id, name) VALUES (2, 'Многодневный пешеходный туризм');
INSERT INTO visit_reasons(id, name) VALUES (3, 'Лыжный туризм');
INSERT INTO visit_reasons(id, name) VALUES (4, 'Спортивные мероприятия');
INSERT INTO visit_reasons(id, name) VALUES (5, 'Научные исследования');
INSERT INTO visit_reasons(id, name) VALUES (6, 'Видео/фотосъемка');
INSERT INTO visit_reasons(id, name) VALUES (7, 'Альпинизм');
INSERT INTO visit_reasons(id, name) VALUES (8, 'Другое');

--feature_types
INSERT INTO feature_types(id, name) VALUES (1, 'входные группы');
INSERT INTO feature_types(id, name) VALUES (2, 'маркировка');
INSERT INTO feature_types(id, name) VALUES (3, 'места для разведения костров');
INSERT INTO feature_types(id, name) VALUES (4, 'навесы для отдыха');
INSERT INTO feature_types(id, name) VALUES (5, 'санитарные зоны');
INSERT INTO feature_types(id, name) VALUES (6, 'кухня-столовая');
INSERT INTO feature_types(id, name) VALUES (7, 'гостевой дом');
INSERT INTO feature_types(id, name) VALUES (8, 'мосты');
INSERT INTO feature_types(id, name) VALUES (9, 'средства навигации и информационные аншлаги');
INSERT INTO feature_types(id, name) VALUES (10, 'дровник');
INSERT INTO feature_types(id, name) VALUES (11, 'визит-центр');

--visit_format
INSERT INTO visit_format(id, name) VALUES (1, 'Многодневный тур(от 1 ночевки и более)');
INSERT INTO visit_format(id, name) VALUES (2, 'Дневная экскурсия (без ночевки на территории парка)');

--permit_statuses
INSERT INTO permit_statuses(id, name) VALUES (1, 'На рассмотрении');
INSERT INTO permit_statuses(id, name) VALUES (2, 'Отказано в посещении');
INSERT INTO permit_statuses(id, name) VALUES (3, 'Посещение одобрено');

--zones
INSERT INTO zones(id, name) VALUES (1, 'Природный парк «Налычево»');

--routes
INSERT INTO routes(id, zone_id, name, description, length_hours, lines) VALUES (1, 1, 'Пиначево-Центральный','1 день: переход 17,5 км. Маршрут начинается от инспекторского кордона и проходит вдоль реки Пиначевская, через стоянку «Промежуточная», до которой от кордона примерно 12 км, к стоянке кордона «Семеновский» - ещё 5,5 км. Ночевка. 2 день: переход 21 км. От кордона «Семеновский», по тропе, примерно в 4 км, находится стоянка «Перевальная». От нее ещё около 5 км, тропа по каменистому распадку и горной тундре уходит вверх к Пиначевскому перевалу, абсолютная высота которого составляет 1150 метров. С перевала открывается великолепная панорама Налычевской долины. Отсюда тропа уходит круто вниз, через кедровый и ольховый стланик. Далее, постепенно выполаживаясь, тропа входит в каменноберезовый лес, и вдоль реки Горячей, спустя 12 км от перевала, подходит к кордону «Центральный». На кордоне, по предварительному бронированию, предусмотрено временное проживание в домиках. В случае использования палаток, туристы размещаются в лагере на р. Жёлтая в 1,5 км от кордона.',48,'{17.5, 21}');

--type_of_reports
INSERT INTO type_of_reports(id, name) VALUES (1, 'Мусор');
INSERT INTO type_of_reports(id, name) VALUES (2, 'Кострища');
INSERT INTO type_of_reports(id, name) VALUES (3, 'Браконьерство');
INSERT INTO type_of_reports(id, name) VALUES (4, 'Пожары');
INSERT INTO type_of_reports(id, name) VALUES (5, 'Другое');

--reports_statuses
INSERT INTO reports_statuses(id, name) VALUES (1, 'В работе');
INSERT INTO reports_statuses(id, name) VALUES (2, 'Отклонено');
INSERT INTO reports_statuses(id, name) VALUES (3, 'Устранено');