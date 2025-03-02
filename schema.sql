--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4
-- Dumped by pg_dump version 17.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: task_progress; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.task_progress (
    task_progress_id integer NOT NULL,
    user_id integer NOT NULL,
    task_id integer NOT NULL,
    status boolean NOT NULL,
    task_date date NOT NULL
);


ALTER TABLE public.task_progress OWNER TO postgres;

--
-- Name: task_progress_progress_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.task_progress_progress_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.task_progress_progress_id_seq OWNER TO postgres;

--
-- Name: task_progress_progress_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.task_progress_progress_id_seq OWNED BY public.task_progress.task_progress_id;


--
-- Name: tasks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tasks (
    task_id integer NOT NULL,
    task_code text NOT NULL,
    task_name text NOT NULL,
    user_ids integer[]
);


ALTER TABLE public.tasks OWNER TO postgres;

--
-- Name: tasks_task_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.tasks_task_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tasks_task_id_seq OWNER TO postgres;

--
-- Name: tasks_task_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tasks_task_id_seq OWNED BY public.tasks.task_id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    user_id integer NOT NULL,
    user_name text NOT NULL,
    pass_word text NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_userid_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_userid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_userid_seq OWNER TO postgres;

--
-- Name: users_userid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_userid_seq OWNED BY public.users.user_id;


--
-- Name: task_progress task_progress_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.task_progress ALTER COLUMN task_progress_id SET DEFAULT nextval('public.task_progress_progress_id_seq'::regclass);


--
-- Name: tasks task_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tasks ALTER COLUMN task_id SET DEFAULT nextval('public.tasks_task_id_seq'::regclass);


--
-- Name: users user_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN user_id SET DEFAULT nextval('public.users_userid_seq'::regclass);


--
-- Data for Name: task_progress; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.task_progress (task_progress_id, user_id, task_id, status, task_date) FROM stdin;
2	2	1	f	2023-10-26
3	3	1	t	2023-10-27
4	2	2	t	2023-10-27
5	4	2	f	2023-10-28
6	1	3	t	2023-10-28
7	5	3	f	2023-10-29
8	3	4	t	2023-10-29
9	4	4	t	2023-10-30
1	1	1	f	2023-10-26
10	1	1	f	2025-03-02
\.


--
-- Data for Name: tasks; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tasks (task_id, task_code, task_name, user_ids) FROM stdin;
4	J1K6L8	Go outside	{3,4}
1	A7B2X9	Go on a walk	{2,3}
9	HE0877	meditate	{1}
10	FPB0KR	new task	{1}
11	7DZ2PV	new task	{1}
12	MF4ORE	new task	{1}
2	3C8D1E	Meditate	{4}
3	F5G9H4	Journal	{5}
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (user_id, user_name, pass_word) FROM stdin;
1	condor25	123
2	jfk	123
3	alice_smith	password123
4	bob_jones	secure456
5	charlie_brown	complex789
6	condor	123
12	jasmine	pw
13	admin	DTeL87qp6CtaejC
14	asnaksdnak	fsafasf
15	asnaksdnak	fsafasf
16	asnaksdnak	dsadsadas
17	jasmine	pe
18	condor	123
19	admin	DTeL87qp6CtaejC
20	sadfs	fdsfa
21	Aj	123
\.


--
-- Name: task_progress_progress_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.task_progress_progress_id_seq', 20, true);


--
-- Name: tasks_task_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.tasks_task_id_seq', 18, true);


--
-- Name: users_userid_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_userid_seq', 21, true);


--
-- Name: task_progress task_progress_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.task_progress
    ADD CONSTRAINT task_progress_pkey PRIMARY KEY (task_progress_id);


--
-- Name: tasks tasks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_pkey PRIMARY KEY (task_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- Name: task_progress task_progress_task_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.task_progress
    ADD CONSTRAINT task_progress_task_id_fkey FOREIGN KEY (task_id) REFERENCES public.tasks(task_id);


--
-- Name: task_progress task_progress_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.task_progress
    ADD CONSTRAINT task_progress_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id);


--
-- PostgreSQL database dump complete
--

