--
-- PostgreSQL database dump
--

\restrict jsdg474bhWte2r1IgbQ75KIWmY4WXI4NAEfaPLtSaDTRIAHhvfEZ3lKRCvCiew5

-- Dumped from database version 18.6 (Debian 18.6-1.pgdg13+2)
-- Dumped by pg_dump version 18.6

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

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: -
--

-- *not* creating schema, since initdb creates it


--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON SCHEMA public IS '';


--
-- Name: postgis; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS postgis WITH SCHEMA public;


--
-- Name: EXTENSION postgis; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION postgis IS 'PostGIS geometry and geography spatial types and functions';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: tlm_aus_einfahrt; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_aus_einfahrt (
    objectid integer NOT NULL,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    nummer character varying(100),
    name character varying(300),
    tlm_strassen_name_uuid character varying,
    shape public.geometry(PointZ,2056)
);


--
-- Name: tlm_aus_einfahrt_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_aus_einfahrt_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_aus_einfahrt_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_aus_einfahrt_objectid_seq OWNED BY public.tlm_aus_einfahrt.objectid;


--
-- Name: tlm_bodenbedeckung; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_bodenbedeckung (
    objectid integer NOT NULL,
    shape_area double precision,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT '-1'::integer,
    shape public.geometry(MultiPolygonZ,2056)
);


--
-- Name: tlm_bodenbedeckung_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_bodenbedeckung_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_bodenbedeckung_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_bodenbedeckung_objectid_seq OWNED BY public.tlm_bodenbedeckung.objectid;


--
-- Name: tlm_einzelbaum; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_einzelbaum (
    objectid integer NOT NULL,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT '-1'::integer,
    shape public.geometry(PointZ,2056)
);


--
-- Name: tlm_einzelbaum_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_einzelbaum_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_einzelbaum_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_einzelbaum_objectid_seq OWNED BY public.tlm_einzelbaum.objectid;


--
-- Name: tlm_einzelobjekt; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_einzelobjekt (
    objectid integer NOT NULL,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 1,
    name character varying(100),
    tlm_eo_name_uuid character varying,
    shape public.geometry(PointZ,2056)
);


--
-- Name: tlm_einzelobjekt_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_einzelobjekt_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_einzelobjekt_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_einzelobjekt_objectid_seq OWNED BY public.tlm_einzelobjekt.objectid;


--
-- Name: tlm_eisenbahn; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_eisenbahn (
    objectid integer NOT NULL,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    kunstbaute integer DEFAULT 100,
    museumsbahn integer DEFAULT 1,
    auf_strasse integer DEFAULT 1,
    achse_dkm integer DEFAULT 1,
    anschlussgleis integer DEFAULT 1,
    eroeffnungsdatum timestamp with time zone,
    stufe integer DEFAULT 0,
    ausser_betrieb integer DEFAULT 1,
    betriebsbahn integer DEFAULT 1,
    standseilbahn integer DEFAULT 1,
    zahnradbahn integer DEFAULT 1,
    anzahl_spuren integer DEFAULT 1,
    verkehrsmittel integer DEFAULT 100,
    name character varying(300),
    tlm_oev_name_uuid character varying,
    shape public.geometry(MultiLineStringZ,2056)
);


--
-- Name: tlm_eisenbahn_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_eisenbahn_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_eisenbahn_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_eisenbahn_objectid_seq OWNED BY public.tlm_eisenbahn.objectid;


--
-- Name: tlm_fliessgewaesser; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_fliessgewaesser (
    objectid integer NOT NULL,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    gewiss_nr integer,
    lauf_nr integer,
    linst character varying(3) DEFAULT 'k_W'::character varying,
    gwl_nr character varying(25) DEFAULT 'k_W'::character varying,
    verlauf integer DEFAULT 100,
    stufe integer DEFAULT 0,
    tlm_gewaesser_lauf_uuid character varying,
    name character varying(300),
    tlm_gewaesser_name_uuid character varying,
    shape public.geometry(MultiLineStringZ,2056)
);


--
-- Name: tlm_fliessgewaesser_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_fliessgewaesser_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_fliessgewaesser_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_fliessgewaesser_objectid_seq OWNED BY public.tlm_fliessgewaesser.objectid;


--
-- Name: tlm_flurname; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_flurname (
    objectid integer NOT NULL,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT '-1'::integer,
    name character varying(300),
    tlm_namen_name_uuid character varying,
    shape public.geometry(PointZ,2056)
);


--
-- Name: tlm_flurname_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_flurname_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_flurname_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_flurname_objectid_seq OWNED BY public.tlm_flurname.objectid;


--
-- Name: tlm_freizeitareal; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_freizeitareal (
    objectid integer NOT NULL,
    shape_area double precision,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    name character varying(300),
    tlm_areale_name_uuid character varying,
    shape public.geometry(MultiPolygonZ,2056)
);


--
-- Name: tlm_freizeitareal_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_freizeitareal_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_freizeitareal_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_freizeitareal_objectid_seq OWNED BY public.tlm_freizeitareal.objectid;


--
-- Name: tlm_gebaeude_footprint; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_gebaeude_footprint (
    objectid integer NOT NULL,
    shape_area double precision,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    nutzung character varying(200),
    tlm_bauten_name_uuid character varying,
    name character varying(300),
    shape public.geometry(MultiPolygonZ,2056)
);


--
-- Name: tlm_gebaeude_footprint_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_gebaeude_footprint_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_gebaeude_footprint_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_gebaeude_footprint_objectid_seq OWNED BY public.tlm_gebaeude_footprint.objectid;


--
-- Name: tlm_gebietsname; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_gebietsname (
    objectid integer NOT NULL,
    shape_area double precision,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    name character varying(300),
    tlm_namen_name_uuid character varying,
    shape public.geometry(MultiPolygonZ,2056)
);


--
-- Name: tlm_gebietsname_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_gebietsname_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_gebietsname_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_gebietsname_objectid_seq OWNED BY public.tlm_gebietsname.objectid;


--
-- Name: tlm_gelaendename; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_gelaendename (
    objectid integer NOT NULL,
    shape_area double precision,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    name character varying(300),
    tlm_namen_name_uuid character varying,
    shape public.geometry(MultiPolygonZ,2056)
);


--
-- Name: tlm_gelaendename_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_gelaendename_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_gelaendename_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_gelaendename_objectid_seq OWNED BY public.tlm_gelaendename.objectid;


--
-- Name: tlm_glamos; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_glamos (
    objectid integer NOT NULL,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer,
    erstellung_monat integer,
    grund_aenderung integer,
    herkunft integer,
    herkunft_jahr integer,
    herkunft_monat integer,
    revision_jahr integer,
    revision_monat integer,
    revision_qualitaet character varying(20),
    tlm_bodenbedeckung_uuid character varying,
    sgi character varying(100),
    riverlevel0 character varying(5),
    riverlevel1 smallint,
    riverlevel2 smallint,
    riverlevel3 character varying(5),
    inventorycode smallint
);


--
-- Name: tlm_glamos_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_glamos_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_glamos_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_glamos_objectid_seq OWNED BY public.tlm_glamos.objectid;


--
-- Name: tlm_haltestelle; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_haltestelle (
    objectid integer NOT NULL,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    dienststellen_nummer integer,
    name character varying(300),
    tlm_oev_name_uuid character varying,
    shape public.geometry(PointZ,2056)
);


--
-- Name: tlm_haltestelle_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_haltestelle_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_haltestelle_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_haltestelle_objectid_seq OWNED BY public.tlm_haltestelle.objectid;


--
-- Name: tlm_leitung; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_leitung (
    objectid integer NOT NULL,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer,
    erstellung_monat integer,
    grund_aenderung integer,
    herkunft integer,
    herkunft_jahr integer,
    herkunft_monat integer,
    revision_jahr integer,
    revision_monat integer,
    revision_qualitaet character varying(20),
    objektart integer DEFAULT 0,
    stufe integer,
    netzebene integer,
    name character varying(100),
    tlm_bauten_name_uuid character varying,
    shape public.geometry(MultiLineStringZ,2056)
);


--
-- Name: tlm_leitung_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_leitung_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_leitung_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_leitung_objectid_seq OWNED BY public.tlm_leitung.objectid;


--
-- Name: tlm_leitung_stromtrasse; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_leitung_stromtrasse (
    rid integer NOT NULL,
    tlm_stromtrasse_uuid character varying,
    tlm_leitung_uuid character varying
);


--
-- Name: tlm_leitung_stromtrasse_rid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_leitung_stromtrasse_rid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_leitung_stromtrasse_rid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_leitung_stromtrasse_rid_seq OWNED BY public.tlm_leitung_stromtrasse.rid;


--
-- Name: tlm_mauer; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_mauer (
    objectid integer NOT NULL,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    shape public.geometry(MultiLineStringZ,2056)
);


--
-- Name: tlm_mauer_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_mauer_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_mauer_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_mauer_objectid_seq OWNED BY public.tlm_mauer.objectid;


--
-- Name: tlm_name_pkt; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_name_pkt (
    objectid integer NOT NULL,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    hoehe double precision DEFAULT 0,
    name character varying(300),
    tlm_namen_name_uuid character varying,
    shape public.geometry(PointZ,2056)
);


--
-- Name: tlm_name_pkt_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_name_pkt_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_name_pkt_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_name_pkt_objectid_seq OWNED BY public.tlm_name_pkt.objectid;


--
-- Name: tlm_nutzungsareal; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_nutzungsareal (
    objectid integer NOT NULL,
    shape_area double precision,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    name character varying(300),
    tlm_areale_name_uuid character varying,
    ris_id character varying(10),
    shape public.geometry(MultiPolygonZ,2056)
);


--
-- Name: tlm_nutzungsareal_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_nutzungsareal_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_nutzungsareal_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_nutzungsareal_objectid_seq OWNED BY public.tlm_nutzungsareal.objectid;


--
-- Name: tlm_nutzungsareal_schule; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_nutzungsareal_schule (
    rid integer NOT NULL,
    tlm_nutzungsareal_uuid character varying,
    tlm_schule_uuid character varying
);


--
-- Name: tlm_nutzungsareal_schule_rid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_nutzungsareal_schule_rid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_nutzungsareal_schule_rid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_nutzungsareal_schule_rid_seq OWNED BY public.tlm_nutzungsareal_schule.rid;


--
-- Name: tlm_schifffahrt; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_schifffahrt (
    objectid integer NOT NULL,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT '-1'::integer,
    name character varying(300),
    tlm_oev_name_uuid character varying,
    shape public.geometry(MultiLineStringZ,2056)
);


--
-- Name: tlm_schifffahrt_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_schifffahrt_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_schifffahrt_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_schifffahrt_objectid_seq OWNED BY public.tlm_schifffahrt.objectid;


--
-- Name: tlm_schule; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_schule (
    objectid integer NOT NULL,
    uuid character varying,
    datum_erstellung timestamp with time zone,
    datum_aenderung timestamp with time zone,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    revision_qualitaet character varying(15),
    isced_stufe integer,
    name_schule character varying(300)
);


--
-- Name: tlm_schule_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_schule_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_schule_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_schule_objectid_seq OWNED BY public.tlm_schule.objectid;


--
-- Name: tlm_schutzgebiet; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_schutzgebiet (
    objectid integer NOT NULL,
    shape_area double precision,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(15),
    objektart integer DEFAULT 0,
    name character varying(300),
    tlm_grenzen_name_uuid character varying,
    shape public.geometry(MultiPolygonZ,2056)
);


--
-- Name: tlm_schutzgebiet_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_schutzgebiet_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_schutzgebiet_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_schutzgebiet_objectid_seq OWNED BY public.tlm_schutzgebiet.objectid;


--
-- Name: tlm_siedlungsname; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_siedlungsname (
    objectid integer NOT NULL,
    shape_area double precision,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    einwohnerkategorie integer DEFAULT 999997,
    name character varying(300),
    tlm_namen_name_uuid character varying,
    shape public.geometry(MultiPolygonZ,2056)
);


--
-- Name: tlm_siedlungsname_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_siedlungsname_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_siedlungsname_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_siedlungsname_objectid_seq OWNED BY public.tlm_siedlungsname.objectid;


--
-- Name: tlm_siedlungsname_zentrum; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_siedlungsname_zentrum (
    objectid integer NOT NULL,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer,
    erstellung_monat integer,
    grund_aenderung integer,
    herkunft integer,
    herkunft_jahr integer,
    herkunft_monat integer,
    revision_jahr integer,
    revision_monat integer,
    revision_qualitaet character varying(20),
    objektart integer DEFAULT 0,
    einwohnerkategorie integer,
    name character varying(300),
    tlm_namen_name_uuid character varying,
    tlm_siedlungsname_uuid character varying,
    shape public.geometry(PointZ,2056)
);


--
-- Name: tlm_siedlungsname_zentrum_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_siedlungsname_zentrum_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_siedlungsname_zentrum_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_siedlungsname_zentrum_objectid_seq OWNED BY public.tlm_siedlungsname_zentrum.objectid;


--
-- Name: tlm_sportbaute_lin; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_sportbaute_lin (
    objectid integer NOT NULL,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    name character varying(300),
    tlm_bauten_name_uuid character varying,
    shape public.geometry(MultiLineStringZ,2056)
);


--
-- Name: tlm_sportbaute_lin_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_sportbaute_lin_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_sportbaute_lin_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_sportbaute_lin_objectid_seq OWNED BY public.tlm_sportbaute_lin.objectid;


--
-- Name: tlm_sportbaute_ply; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_sportbaute_ply (
    objectid integer NOT NULL,
    shape_area double precision,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    name character varying(300),
    tlm_bauten_name_uuid character varying,
    shape public.geometry(MultiPolygonZ,2056)
);


--
-- Name: tlm_sportbaute_ply_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_sportbaute_ply_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_sportbaute_ply_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_sportbaute_ply_objectid_seq OWNED BY public.tlm_sportbaute_ply.objectid;


--
-- Name: tlm_staubaute; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_staubaute (
    objectid integer NOT NULL,
    shape_area double precision,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    name character varying(300),
    tlm_bauten_name_uuid character varying,
    shape public.geometry(MultiPolygonZ,2056)
);


--
-- Name: tlm_staubaute_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_staubaute_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_staubaute_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_staubaute_objectid_seq OWNED BY public.tlm_staubaute.objectid;


--
-- Name: tlm_stehendes_gewaesser; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_stehendes_gewaesser (
    objectid integer NOT NULL,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_monat integer DEFAULT 1,
    erstellung_jahr integer DEFAULT 1900,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    gewiss_nr integer,
    lauf_nr integer,
    linst character varying(3) DEFAULT 'k_W'::character varying,
    gwl_nr character varying(25) DEFAULT 'k_W'::character varying,
    tlm_gewaesser_lauf_uuid character varying,
    name character varying(300),
    tlm_gewaesser_name_uuid character varying,
    wasserstand_wechselnd integer,
    shape public.geometry(MultiLineStringZ,2056)
);


--
-- Name: tlm_stehendes_gewaesser_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_stehendes_gewaesser_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_stehendes_gewaesser_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_stehendes_gewaesser_objectid_seq OWNED BY public.tlm_stehendes_gewaesser.objectid;


--
-- Name: tlm_strasse; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_strasse (
    objectid integer NOT NULL,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    kunstbaute integer DEFAULT 100,
    wanderwege integer,
    verkehrsbeschraenkung integer,
    befahrbarkeit integer,
    eroeffnungsdatum timestamp with time zone,
    stufe integer DEFAULT 0,
    richtungsgetrennt integer DEFAULT 1,
    belagsart integer DEFAULT 100,
    kreisel integer DEFAULT 1,
    eigentuemer integer,
    verkehrsbedeutung integer,
    name character varying(300),
    tlm_strassen_name_uuid character varying,
    strassenname character varying(150),
    shape public.geometry(MultiLineStringZ,2056)
);


--
-- Name: tlm_strasse_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_strasse_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_strasse_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_strasse_objectid_seq OWNED BY public.tlm_strasse.objectid;


--
-- Name: tlm_strasseninfo; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_strasseninfo (
    objectid integer NOT NULL,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    tankstelle integer DEFAULT 1,
    restaurant integer DEFAULT 1,
    shop integer DEFAULT 1,
    toiletten integer DEFAULT 1,
    name character varying(300),
    tlm_strassen_name_uuid character varying,
    shape public.geometry(PointZ,2056)
);


--
-- Name: tlm_strasseninfo_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_strasseninfo_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_strasseninfo_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_strasseninfo_objectid_seq OWNED BY public.tlm_strasseninfo.objectid;


--
-- Name: tlm_strassenname; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_strassenname (
    objectid integer NOT NULL,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer,
    erstellung_monat integer,
    grund_aenderung integer,
    herkunft integer,
    herkunft_jahr integer,
    herkunft_monat integer,
    revision_jahr integer,
    revision_monat integer,
    revision_qualitaet character varying(20),
    name character varying(60),
    esid double precision,
    kanton character varying(2),
    gdename character varying(40),
    gdenr smallint,
    sprachcode character varying(3)
);


--
-- Name: tlm_strassenname_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_strassenname_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_strassenname_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_strassenname_objectid_seq OWNED BY public.tlm_strassenname.objectid;


--
-- Name: tlm_strassenname_strasse; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_strassenname_strasse (
    rid integer NOT NULL,
    tlm_strassenname_uuid character varying,
    tlm_strasse_uuid character varying
);


--
-- Name: tlm_strassenname_strasse_rid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_strassenname_strasse_rid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_strassenname_strasse_rid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_strassenname_strasse_rid_seq OWNED BY public.tlm_strassenname_strasse.rid;


--
-- Name: tlm_strassenroute; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_strassenroute (
    objectid integer NOT NULL,
    uuid character varying,
    datum_erstellung timestamp with time zone,
    datum_aenderung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(15),
    objektart integer DEFAULT 0,
    name character varying(300),
    routennummer character varying(10)
);


--
-- Name: tlm_strassenroute_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_strassenroute_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_strassenroute_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_strassenroute_objectid_seq OWNED BY public.tlm_strassenroute.objectid;


--
-- Name: tlm_strassenroute_strasse; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_strassenroute_strasse (
    rid integer NOT NULL,
    tlm_strassenroute_uuid character varying,
    tlm_strasse_uuid character varying
);


--
-- Name: tlm_strassenroute_strasse_rid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_strassenroute_strasse_rid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_strassenroute_strasse_rid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_strassenroute_strasse_rid_seq OWNED BY public.tlm_strassenroute_strasse.rid;


--
-- Name: tlm_stromtrasse; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_stromtrasse (
    objectid integer NOT NULL,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer,
    erstellung_monat integer,
    grund_aenderung integer,
    herkunft integer,
    herkunft_jahr integer,
    herkunft_monat integer,
    revision_jahr integer,
    revision_monat integer,
    revision_qualitaet character varying(20),
    netzebene integer
);


--
-- Name: tlm_stromtrasse_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_stromtrasse_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_stromtrasse_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_stromtrasse_objectid_seq OWNED BY public.tlm_stromtrasse.objectid;


--
-- Name: tlm_uebrige_bahn; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_uebrige_bahn (
    objectid integer NOT NULL,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    betriebsbahn integer,
    stufe integer DEFAULT 0,
    ausser_betrieb integer DEFAULT 1,
    eroeffnungsdatum timestamp with time zone,
    name character varying(300),
    tlm_oev_name_uuid character varying,
    shape public.geometry(MultiLineStringZ,2056)
);


--
-- Name: tlm_uebrige_bahn_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_uebrige_bahn_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_uebrige_bahn_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_uebrige_bahn_objectid_seq OWNED BY public.tlm_uebrige_bahn.objectid;


--
-- Name: tlm_verbauung; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_verbauung (
    objectid integer NOT NULL,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    shape public.geometry(MultiLineStringZ,2056)
);


--
-- Name: tlm_verbauung_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_verbauung_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_verbauung_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_verbauung_objectid_seq OWNED BY public.tlm_verbauung.objectid;


--
-- Name: tlm_verkehrsareal; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_verkehrsareal (
    objectid integer NOT NULL,
    shape_area double precision,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 0,
    name character varying(300),
    tlm_areale_name_uuid character varying,
    shape public.geometry(MultiPolygonZ,2056)
);


--
-- Name: tlm_verkehrsareal_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_verkehrsareal_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_verkehrsareal_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_verkehrsareal_objectid_seq OWNED BY public.tlm_verkehrsareal.objectid;


--
-- Name: tlm_verkehrsbaute_lin; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_verkehrsbaute_lin (
    objectid integer NOT NULL,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    objektart integer DEFAULT 0,
    stufe integer DEFAULT 0,
    revision_qualitaet character varying(15),
    shape public.geometry(MultiLineStringZ,2056)
);


--
-- Name: tlm_verkehrsbaute_lin_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_verkehrsbaute_lin_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_verkehrsbaute_lin_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_verkehrsbaute_lin_objectid_seq OWNED BY public.tlm_verkehrsbaute_lin.objectid;


--
-- Name: tlm_verkehrsbaute_ply; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_verkehrsbaute_ply (
    objectid integer NOT NULL,
    shape_area double precision,
    shape_length double precision,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT '-1'::integer,
    shape public.geometry(MultiPolygonZ,2056)
);


--
-- Name: tlm_verkehrsbaute_ply_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_verkehrsbaute_ply_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_verkehrsbaute_ply_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_verkehrsbaute_ply_objectid_seq OWNED BY public.tlm_verkehrsbaute_ply.objectid;


--
-- Name: tlm_versorgungsbaute_pkt; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tlm_versorgungsbaute_pkt (
    objectid integer NOT NULL,
    uuid character varying,
    datum_aenderung timestamp with time zone,
    datum_erstellung timestamp with time zone,
    erstellung_jahr integer DEFAULT 1900,
    erstellung_monat integer DEFAULT 1,
    revision_jahr integer DEFAULT 1900,
    revision_monat integer DEFAULT 1,
    grund_aenderung integer DEFAULT 200,
    herkunft integer DEFAULT 100,
    herkunft_jahr integer DEFAULT 1900,
    herkunft_monat integer DEFAULT 1,
    revision_qualitaet character varying(100),
    objektart integer DEFAULT 2,
    shape public.geometry(PointZ,2056)
);


--
-- Name: tlm_versorgungsbaute_pkt_objectid_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tlm_versorgungsbaute_pkt_objectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tlm_versorgungsbaute_pkt_objectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tlm_versorgungsbaute_pkt_objectid_seq OWNED BY public.tlm_versorgungsbaute_pkt.objectid;


--
-- Name: tlm_aus_einfahrt objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_aus_einfahrt ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_aus_einfahrt_objectid_seq'::regclass);


--
-- Name: tlm_bodenbedeckung objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_bodenbedeckung ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_bodenbedeckung_objectid_seq'::regclass);


--
-- Name: tlm_einzelbaum objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_einzelbaum ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_einzelbaum_objectid_seq'::regclass);


--
-- Name: tlm_einzelobjekt objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_einzelobjekt ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_einzelobjekt_objectid_seq'::regclass);


--
-- Name: tlm_eisenbahn objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_eisenbahn ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_eisenbahn_objectid_seq'::regclass);


--
-- Name: tlm_fliessgewaesser objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_fliessgewaesser ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_fliessgewaesser_objectid_seq'::regclass);


--
-- Name: tlm_flurname objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_flurname ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_flurname_objectid_seq'::regclass);


--
-- Name: tlm_freizeitareal objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_freizeitareal ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_freizeitareal_objectid_seq'::regclass);


--
-- Name: tlm_gebaeude_footprint objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_gebaeude_footprint ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_gebaeude_footprint_objectid_seq'::regclass);


--
-- Name: tlm_gebietsname objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_gebietsname ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_gebietsname_objectid_seq'::regclass);


--
-- Name: tlm_gelaendename objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_gelaendename ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_gelaendename_objectid_seq'::regclass);


--
-- Name: tlm_glamos objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_glamos ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_glamos_objectid_seq'::regclass);


--
-- Name: tlm_haltestelle objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_haltestelle ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_haltestelle_objectid_seq'::regclass);


--
-- Name: tlm_leitung objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_leitung ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_leitung_objectid_seq'::regclass);


--
-- Name: tlm_leitung_stromtrasse rid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_leitung_stromtrasse ALTER COLUMN rid SET DEFAULT nextval('public.tlm_leitung_stromtrasse_rid_seq'::regclass);


--
-- Name: tlm_mauer objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_mauer ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_mauer_objectid_seq'::regclass);


--
-- Name: tlm_name_pkt objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_name_pkt ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_name_pkt_objectid_seq'::regclass);


--
-- Name: tlm_nutzungsareal objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_nutzungsareal ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_nutzungsareal_objectid_seq'::regclass);


--
-- Name: tlm_nutzungsareal_schule rid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_nutzungsareal_schule ALTER COLUMN rid SET DEFAULT nextval('public.tlm_nutzungsareal_schule_rid_seq'::regclass);


--
-- Name: tlm_schifffahrt objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_schifffahrt ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_schifffahrt_objectid_seq'::regclass);


--
-- Name: tlm_schule objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_schule ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_schule_objectid_seq'::regclass);


--
-- Name: tlm_schutzgebiet objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_schutzgebiet ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_schutzgebiet_objectid_seq'::regclass);


--
-- Name: tlm_siedlungsname objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_siedlungsname ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_siedlungsname_objectid_seq'::regclass);


--
-- Name: tlm_siedlungsname_zentrum objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_siedlungsname_zentrum ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_siedlungsname_zentrum_objectid_seq'::regclass);


--
-- Name: tlm_sportbaute_lin objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_sportbaute_lin ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_sportbaute_lin_objectid_seq'::regclass);


--
-- Name: tlm_sportbaute_ply objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_sportbaute_ply ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_sportbaute_ply_objectid_seq'::regclass);


--
-- Name: tlm_staubaute objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_staubaute ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_staubaute_objectid_seq'::regclass);


--
-- Name: tlm_stehendes_gewaesser objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_stehendes_gewaesser ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_stehendes_gewaesser_objectid_seq'::regclass);


--
-- Name: tlm_strasse objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_strasse ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_strasse_objectid_seq'::regclass);


--
-- Name: tlm_strasseninfo objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_strasseninfo ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_strasseninfo_objectid_seq'::regclass);


--
-- Name: tlm_strassenname objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_strassenname ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_strassenname_objectid_seq'::regclass);


--
-- Name: tlm_strassenname_strasse rid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_strassenname_strasse ALTER COLUMN rid SET DEFAULT nextval('public.tlm_strassenname_strasse_rid_seq'::regclass);


--
-- Name: tlm_strassenroute objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_strassenroute ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_strassenroute_objectid_seq'::regclass);


--
-- Name: tlm_strassenroute_strasse rid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_strassenroute_strasse ALTER COLUMN rid SET DEFAULT nextval('public.tlm_strassenroute_strasse_rid_seq'::regclass);


--
-- Name: tlm_stromtrasse objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_stromtrasse ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_stromtrasse_objectid_seq'::regclass);


--
-- Name: tlm_uebrige_bahn objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_uebrige_bahn ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_uebrige_bahn_objectid_seq'::regclass);


--
-- Name: tlm_verbauung objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_verbauung ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_verbauung_objectid_seq'::regclass);


--
-- Name: tlm_verkehrsareal objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_verkehrsareal ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_verkehrsareal_objectid_seq'::regclass);


--
-- Name: tlm_verkehrsbaute_lin objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_verkehrsbaute_lin ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_verkehrsbaute_lin_objectid_seq'::regclass);


--
-- Name: tlm_verkehrsbaute_ply objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_verkehrsbaute_ply ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_verkehrsbaute_ply_objectid_seq'::regclass);


--
-- Name: tlm_versorgungsbaute_pkt objectid; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_versorgungsbaute_pkt ALTER COLUMN objectid SET DEFAULT nextval('public.tlm_versorgungsbaute_pkt_objectid_seq'::regclass);


--
-- Name: tlm_aus_einfahrt tlm_aus_einfahrt_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_aus_einfahrt
    ADD CONSTRAINT tlm_aus_einfahrt_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_bodenbedeckung tlm_bodenbedeckung_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_bodenbedeckung
    ADD CONSTRAINT tlm_bodenbedeckung_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_einzelbaum tlm_einzelbaum_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_einzelbaum
    ADD CONSTRAINT tlm_einzelbaum_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_einzelobjekt tlm_einzelobjekt_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_einzelobjekt
    ADD CONSTRAINT tlm_einzelobjekt_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_eisenbahn tlm_eisenbahn_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_eisenbahn
    ADD CONSTRAINT tlm_eisenbahn_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_fliessgewaesser tlm_fliessgewaesser_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_fliessgewaesser
    ADD CONSTRAINT tlm_fliessgewaesser_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_flurname tlm_flurname_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_flurname
    ADD CONSTRAINT tlm_flurname_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_freizeitareal tlm_freizeitareal_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_freizeitareal
    ADD CONSTRAINT tlm_freizeitareal_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_gebaeude_footprint tlm_gebaeude_footprint_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_gebaeude_footprint
    ADD CONSTRAINT tlm_gebaeude_footprint_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_gebietsname tlm_gebietsname_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_gebietsname
    ADD CONSTRAINT tlm_gebietsname_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_gelaendename tlm_gelaendename_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_gelaendename
    ADD CONSTRAINT tlm_gelaendename_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_glamos tlm_glamos_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_glamos
    ADD CONSTRAINT tlm_glamos_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_haltestelle tlm_haltestelle_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_haltestelle
    ADD CONSTRAINT tlm_haltestelle_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_leitung tlm_leitung_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_leitung
    ADD CONSTRAINT tlm_leitung_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_leitung_stromtrasse tlm_leitung_stromtrasse_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_leitung_stromtrasse
    ADD CONSTRAINT tlm_leitung_stromtrasse_pkey PRIMARY KEY (rid);


--
-- Name: tlm_mauer tlm_mauer_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_mauer
    ADD CONSTRAINT tlm_mauer_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_name_pkt tlm_name_pkt_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_name_pkt
    ADD CONSTRAINT tlm_name_pkt_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_nutzungsareal tlm_nutzungsareal_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_nutzungsareal
    ADD CONSTRAINT tlm_nutzungsareal_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_nutzungsareal_schule tlm_nutzungsareal_schule_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_nutzungsareal_schule
    ADD CONSTRAINT tlm_nutzungsareal_schule_pkey PRIMARY KEY (rid);


--
-- Name: tlm_schifffahrt tlm_schifffahrt_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_schifffahrt
    ADD CONSTRAINT tlm_schifffahrt_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_schule tlm_schule_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_schule
    ADD CONSTRAINT tlm_schule_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_schutzgebiet tlm_schutzgebiet_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_schutzgebiet
    ADD CONSTRAINT tlm_schutzgebiet_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_siedlungsname tlm_siedlungsname_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_siedlungsname
    ADD CONSTRAINT tlm_siedlungsname_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_siedlungsname_zentrum tlm_siedlungsname_zentrum_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_siedlungsname_zentrum
    ADD CONSTRAINT tlm_siedlungsname_zentrum_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_sportbaute_lin tlm_sportbaute_lin_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_sportbaute_lin
    ADD CONSTRAINT tlm_sportbaute_lin_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_sportbaute_ply tlm_sportbaute_ply_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_sportbaute_ply
    ADD CONSTRAINT tlm_sportbaute_ply_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_staubaute tlm_staubaute_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_staubaute
    ADD CONSTRAINT tlm_staubaute_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_stehendes_gewaesser tlm_stehendes_gewaesser_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_stehendes_gewaesser
    ADD CONSTRAINT tlm_stehendes_gewaesser_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_strasse tlm_strasse_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_strasse
    ADD CONSTRAINT tlm_strasse_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_strasseninfo tlm_strasseninfo_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_strasseninfo
    ADD CONSTRAINT tlm_strasseninfo_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_strassenname tlm_strassenname_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_strassenname
    ADD CONSTRAINT tlm_strassenname_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_strassenname_strasse tlm_strassenname_strasse_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_strassenname_strasse
    ADD CONSTRAINT tlm_strassenname_strasse_pkey PRIMARY KEY (rid);


--
-- Name: tlm_strassenroute tlm_strassenroute_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_strassenroute
    ADD CONSTRAINT tlm_strassenroute_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_strassenroute_strasse tlm_strassenroute_strasse_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_strassenroute_strasse
    ADD CONSTRAINT tlm_strassenroute_strasse_pkey PRIMARY KEY (rid);


--
-- Name: tlm_stromtrasse tlm_stromtrasse_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_stromtrasse
    ADD CONSTRAINT tlm_stromtrasse_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_uebrige_bahn tlm_uebrige_bahn_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_uebrige_bahn
    ADD CONSTRAINT tlm_uebrige_bahn_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_verbauung tlm_verbauung_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_verbauung
    ADD CONSTRAINT tlm_verbauung_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_verkehrsareal tlm_verkehrsareal_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_verkehrsareal
    ADD CONSTRAINT tlm_verkehrsareal_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_verkehrsbaute_lin tlm_verkehrsbaute_lin_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_verkehrsbaute_lin
    ADD CONSTRAINT tlm_verkehrsbaute_lin_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_verkehrsbaute_ply tlm_verkehrsbaute_ply_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_verkehrsbaute_ply
    ADD CONSTRAINT tlm_verkehrsbaute_ply_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_versorgungsbaute_pkt tlm_versorgungsbaute_pkt_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tlm_versorgungsbaute_pkt
    ADD CONSTRAINT tlm_versorgungsbaute_pkt_pkey PRIMARY KEY (objectid);


--
-- Name: tlm_aus_einfahrt_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_aus_einfahrt_shape_geom_idx ON public.tlm_aus_einfahrt USING gist (shape);


--
-- Name: tlm_bodenbedeckung_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_bodenbedeckung_shape_geom_idx ON public.tlm_bodenbedeckung USING gist (shape);


--
-- Name: tlm_einzelbaum_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_einzelbaum_shape_geom_idx ON public.tlm_einzelbaum USING gist (shape);


--
-- Name: tlm_einzelobjekt_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_einzelobjekt_shape_geom_idx ON public.tlm_einzelobjekt USING gist (shape);


--
-- Name: tlm_eisenbahn_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_eisenbahn_shape_geom_idx ON public.tlm_eisenbahn USING gist (shape);


--
-- Name: tlm_fliessgewaesser_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_fliessgewaesser_shape_geom_idx ON public.tlm_fliessgewaesser USING gist (shape);


--
-- Name: tlm_flurname_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_flurname_shape_geom_idx ON public.tlm_flurname USING gist (shape);


--
-- Name: tlm_freizeitareal_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_freizeitareal_shape_geom_idx ON public.tlm_freizeitareal USING gist (shape);


--
-- Name: tlm_gebaeude_footprint_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_gebaeude_footprint_shape_geom_idx ON public.tlm_gebaeude_footprint USING gist (shape);


--
-- Name: tlm_gebietsname_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_gebietsname_shape_geom_idx ON public.tlm_gebietsname USING gist (shape);


--
-- Name: tlm_gelaendename_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_gelaendename_shape_geom_idx ON public.tlm_gelaendename USING gist (shape);


--
-- Name: tlm_haltestelle_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_haltestelle_shape_geom_idx ON public.tlm_haltestelle USING gist (shape);


--
-- Name: tlm_leitung_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_leitung_shape_geom_idx ON public.tlm_leitung USING gist (shape);


--
-- Name: tlm_mauer_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_mauer_shape_geom_idx ON public.tlm_mauer USING gist (shape);


--
-- Name: tlm_name_pkt_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_name_pkt_shape_geom_idx ON public.tlm_name_pkt USING gist (shape);


--
-- Name: tlm_nutzungsareal_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_nutzungsareal_shape_geom_idx ON public.tlm_nutzungsareal USING gist (shape);


--
-- Name: tlm_schifffahrt_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_schifffahrt_shape_geom_idx ON public.tlm_schifffahrt USING gist (shape);


--
-- Name: tlm_schutzgebiet_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_schutzgebiet_shape_geom_idx ON public.tlm_schutzgebiet USING gist (shape);


--
-- Name: tlm_siedlungsname_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_siedlungsname_shape_geom_idx ON public.tlm_siedlungsname USING gist (shape);


--
-- Name: tlm_siedlungsname_zentrum_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_siedlungsname_zentrum_shape_geom_idx ON public.tlm_siedlungsname_zentrum USING gist (shape);


--
-- Name: tlm_sportbaute_lin_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_sportbaute_lin_shape_geom_idx ON public.tlm_sportbaute_lin USING gist (shape);


--
-- Name: tlm_sportbaute_ply_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_sportbaute_ply_shape_geom_idx ON public.tlm_sportbaute_ply USING gist (shape);


--
-- Name: tlm_staubaute_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_staubaute_shape_geom_idx ON public.tlm_staubaute USING gist (shape);


--
-- Name: tlm_stehendes_gewaesser_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_stehendes_gewaesser_shape_geom_idx ON public.tlm_stehendes_gewaesser USING gist (shape);


--
-- Name: tlm_strasse_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_strasse_shape_geom_idx ON public.tlm_strasse USING gist (shape);


--
-- Name: tlm_strasseninfo_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_strasseninfo_shape_geom_idx ON public.tlm_strasseninfo USING gist (shape);


--
-- Name: tlm_uebrige_bahn_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_uebrige_bahn_shape_geom_idx ON public.tlm_uebrige_bahn USING gist (shape);


--
-- Name: tlm_verbauung_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_verbauung_shape_geom_idx ON public.tlm_verbauung USING gist (shape);


--
-- Name: tlm_verkehrsareal_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_verkehrsareal_shape_geom_idx ON public.tlm_verkehrsareal USING gist (shape);


--
-- Name: tlm_verkehrsbaute_lin_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_verkehrsbaute_lin_shape_geom_idx ON public.tlm_verkehrsbaute_lin USING gist (shape);


--
-- Name: tlm_verkehrsbaute_ply_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_verkehrsbaute_ply_shape_geom_idx ON public.tlm_verkehrsbaute_ply USING gist (shape);


--
-- Name: tlm_versorgungsbaute_pkt_shape_geom_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tlm_versorgungsbaute_pkt_shape_geom_idx ON public.tlm_versorgungsbaute_pkt USING gist (shape);


--
-- PostgreSQL database dump complete
--

\unrestrict jsdg474bhWte2r1IgbQ75KIWmY4WXI4NAEfaPLtSaDTRIAHhvfEZ3lKRCvCiew5

