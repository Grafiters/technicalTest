--
-- PostgreSQL database dump
--

-- Dumped from database version 13.21
-- Dumped by pg_dump version 13.21

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
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
-- Name: customers; Type: TABLE; Schema: public; Owner: Admin
--

CREATE TABLE public.customers (
    id bigint NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255),
    nik integer,
    full_name character varying(255),
    legal_name character varying(255),
    birth_place character varying(255),
    birth_date date,
    salary bigint,
    ktp_image_url text,
    selfie_image_url text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.customers OWNER TO "Admin";

--
-- Name: customers_id_seq; Type: SEQUENCE; Schema: public; Owner: Admin
--

CREATE SEQUENCE public.customers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.customers_id_seq OWNER TO "Admin";

--
-- Name: customers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: Admin
--

ALTER SEQUENCE public.customers_id_seq OWNED BY public.customers.id;


--
-- Name: installment_logs; Type: TABLE; Schema: public; Owner: Admin
--

CREATE TABLE public.installment_logs (
    id bigint NOT NULL,
    transaction_id bigint NOT NULL,
    month integer,
    amount bigint,
    due_date date,
    paid_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.installment_logs OWNER TO "Admin";

--
-- Name: installment_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: Admin
--

CREATE SEQUENCE public.installment_logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.installment_logs_id_seq OWNER TO "Admin";

--
-- Name: installment_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: Admin
--

ALTER SEQUENCE public.installment_logs_id_seq OWNED BY public.installment_logs.id;


--
-- Name: limits; Type: TABLE; Schema: public; Owner: Admin
--

CREATE TABLE public.limits (
    id bigint NOT NULL,
    customer_id bigint NOT NULL,
    tenor integer,
    amount bigint,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.limits OWNER TO "Admin";

--
-- Name: limits_id_seq; Type: SEQUENCE; Schema: public; Owner: Admin
--

CREATE SEQUENCE public.limits_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.limits_id_seq OWNER TO "Admin";

--
-- Name: limits_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: Admin
--

ALTER SEQUENCE public.limits_id_seq OWNED BY public.limits.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: Admin
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean
);


ALTER TABLE public.schema_migrations OWNER TO "Admin";

--
-- Name: transactions; Type: TABLE; Schema: public; Owner: Admin
--

CREATE TABLE public.transactions (
    id bigint NOT NULL,
    customer_id bigint NOT NULL,
    limit_id bigint NOT NULL,
    contract_no bigint,
    otr bigint,
    admin_fee bigint,
    installment bigint,
    asset_name character varying(255),
    status character varying(50),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.transactions OWNER TO "Admin";

--
-- Name: transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: Admin
--

CREATE SEQUENCE public.transactions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.transactions_id_seq OWNER TO "Admin";

--
-- Name: transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: Admin
--

ALTER SEQUENCE public.transactions_id_seq OWNED BY public.transactions.id;


--
-- Name: customers id; Type: DEFAULT; Schema: public; Owner: Admin
--

ALTER TABLE ONLY public.customers ALTER COLUMN id SET DEFAULT nextval('public.customers_id_seq'::regclass);


--
-- Name: installment_logs id; Type: DEFAULT; Schema: public; Owner: Admin
--

ALTER TABLE ONLY public.installment_logs ALTER COLUMN id SET DEFAULT nextval('public.installment_logs_id_seq'::regclass);


--
-- Name: limits id; Type: DEFAULT; Schema: public; Owner: Admin
--

ALTER TABLE ONLY public.limits ALTER COLUMN id SET DEFAULT nextval('public.limits_id_seq'::regclass);


--
-- Name: transactions id; Type: DEFAULT; Schema: public; Owner: Admin
--

ALTER TABLE ONLY public.transactions ALTER COLUMN id SET DEFAULT nextval('public.transactions_id_seq'::regclass);


--
-- Data for Name: customers; Type: TABLE DATA; Schema: public; Owner: Admin
--

COPY public.customers (id, email, password, nik, full_name, legal_name, birth_place, birth_date, salary, ktp_image_url, selfie_image_url, created_at, updated_at) FROM stdin;
1	alone@gmail.com	$2a$10$rhWn2yqQEZw2EXoA2E3qZOj7NyueRdUC9/gkEXfdg4/t7pxMlzic.	0	alone	alone	grobogan	2000-08-04	10000000	customer/1/079d5d8238fe45afbc6873f7e291cf3e4beef5fc.png	customer/1/079d5d8238fe45afbc6873f7e291cf3e4beef5fc.png	2025-07-28 19:06:06.634169	2025-07-29 11:18:44.940039
\.


--
-- Data for Name: installment_logs; Type: TABLE DATA; Schema: public; Owner: Admin
--

COPY public.installment_logs (id, transaction_id, month, amount, due_date, paid_at, created_at, updated_at) FROM stdin;
4	10	3	625000	2025-10-01	\N	2025-07-29 00:10:10.707111	2025-07-29 00:10:10.707111
5	10	4	625000	2025-11-01	\N	2025-07-29 00:10:10.711999	2025-07-29 00:10:10.711999
6	10	5	625000	2025-12-01	\N	2025-07-29 00:10:10.718412	2025-07-29 00:10:10.718412
7	10	6	625000	2026-01-01	\N	2025-07-29 00:10:10.723995	2025-07-29 00:10:10.723995
2	10	1	625000	2025-08-01	2025-07-29 00:00:00	2025-07-29 00:10:10.695815	2025-07-29 01:41:28.691064
3	10	2	625000	2025-09-01	2025-07-29 00:00:00	2025-07-29 00:10:10.701051	2025-07-29 01:42:56.935725
\.


--
-- Data for Name: limits; Type: TABLE DATA; Schema: public; Owner: Admin
--

COPY public.limits (id, customer_id, tenor, amount, created_at, updated_at) FROM stdin;
5	1	2	6000000	2025-07-28 19:47:39.558896	2025-07-28 19:47:39.558896
6	1	3	7500000	2025-07-28 19:47:39.558896	2025-07-28 19:47:39.558896
8	1	1	4000000	2025-07-28 19:47:39.558896	2025-07-28 19:47:39.558896
7	1	6	0	2025-07-28 19:47:39.558896	2025-07-28 19:47:39.558896
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: Admin
--

COPY public.schema_migrations (version, dirty) FROM stdin;
20250728112417	\N
20250728114728	\N
20250728115001	\N
20250728115034	\N
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: Admin
--

COPY public.transactions (id, customer_id, limit_id, contract_no, otr, admin_fee, installment, asset_name, status, created_at, updated_at) FROM stdin;
10	1	7	123456789	4750000	1500000	3750000	Toyota Avanza 1.5 G AT 2022	active	2025-07-29 00:10:10.682813	2025-07-29 01:18:53.425703
\.


--
-- Name: customers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: Admin
--

SELECT pg_catalog.setval('public.customers_id_seq', 4, true);


--
-- Name: installment_logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: Admin
--

SELECT pg_catalog.setval('public.installment_logs_id_seq', 7, true);


--
-- Name: limits_id_seq; Type: SEQUENCE SET; Schema: public; Owner: Admin
--

SELECT pg_catalog.setval('public.limits_id_seq', 8, true);


--
-- Name: transactions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: Admin
--

SELECT pg_catalog.setval('public.transactions_id_seq', 10, true);


--
-- Name: customers customers_email_key; Type: CONSTRAINT; Schema: public; Owner: Admin
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_email_key UNIQUE (email);


--
-- Name: customers customers_pkey; Type: CONSTRAINT; Schema: public; Owner: Admin
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY (id);


--
-- Name: installment_logs installment_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: Admin
--

ALTER TABLE ONLY public.installment_logs
    ADD CONSTRAINT installment_logs_pkey PRIMARY KEY (id);


--
-- Name: limits limits_pkey; Type: CONSTRAINT; Schema: public; Owner: Admin
--

ALTER TABLE ONLY public.limits
    ADD CONSTRAINT limits_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: Admin
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: Admin
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: idx_email; Type: INDEX; Schema: public; Owner: Admin
--

CREATE INDEX idx_email ON public.customers USING btree (email);


--
-- Name: idx_nik; Type: INDEX; Schema: public; Owner: Admin
--

CREATE INDEX idx_nik ON public.customers USING btree (nik);


--
-- Name: installment_logs installment_logs_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: Admin
--

ALTER TABLE ONLY public.installment_logs
    ADD CONSTRAINT installment_logs_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES public.transactions(id) ON DELETE CASCADE;


--
-- Name: limits limits_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: Admin
--

ALTER TABLE ONLY public.limits
    ADD CONSTRAINT limits_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE CASCADE;


--
-- Name: transactions transactions_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: Admin
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE CASCADE;


--
-- Name: transactions transactions_limit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: Admin
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_limit_id_fkey FOREIGN KEY (limit_id) REFERENCES public.limits(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

