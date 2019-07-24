--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.5
-- Dumped by pg_dump version 10.5 (Ubuntu 10.5-1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: activities; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.activities (
    id uuid NOT NULL,
    rpc text NOT NULL,
    data jsonb NOT NULL,
    error text,
    user_id uuid,
    created_at timestamp without time zone NOT NULL
);


--
-- Name: announcements; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.announcements (
    id uuid NOT NULL,
    message text DEFAULT ''::text NOT NULL,
    card_id uuid,
    from_user uuid NOT NULL,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


--
-- Name: anonymous_aliases; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.anonymous_aliases (
    id uuid NOT NULL,
    display_name text NOT NULL,
    username text NOT NULL,
    profile_image_path text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    inactive boolean DEFAULT false NOT NULL,
    CONSTRAINT anonymous_aliases_display_name_check CHECK ((display_name <> ''::text)),
    CONSTRAINT anonymous_aliases_profile_image_path_check CHECK ((profile_image_path <> ''::text)),
    CONSTRAINT anonymous_aliases_username_check CHECK ((username <> ''::text))
);


--
-- Name: cards; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cards (
    id uuid NOT NULL,
    owner_id uuid NOT NULL,
    thread_reply_id uuid,
    thread_root_id uuid,
    title text NOT NULL,
    content text NOT NULL,
    url text NOT NULL,
    anonymous boolean NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    background_image_path character varying(255),
    background_color character varying(255),
    author_to_alias jsonb DEFAULT '{}'::jsonb NOT NULL,
    alias_id uuid,
    deleted_at timestamp without time zone,
    shadowbanned_at timestamp without time zone,
    channel_id uuid,
    is_intro_card boolean DEFAULT false NOT NULL,
    coins_earned integer DEFAULT 0 NOT NULL,
    CONSTRAINT cards_thread_root_id_check CHECK (((thread_reply_id IS NULL) OR (thread_root_id IS NOT NULL)))
);


--
-- Name: channel_memberships; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.channel_memberships (
    user_id uuid NOT NULL,
    channel_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: channel_mutes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.channel_mutes (
    user_id uuid NOT NULL,
    channel_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: channels; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.channels (
    id uuid NOT NULL,
    handle text NOT NULL,
    name text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    is_default boolean DEFAULT false NOT NULL,
    owner_id uuid,
    description text DEFAULT ''::text NOT NULL,
    private boolean DEFAULT false NOT NULL
);


--
-- Name: coin_rewards; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.coin_rewards (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    coins_received integer DEFAULT 0 NOT NULL,
    last_rewarded_on timestamp without time zone,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: coin_transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.coin_transactions (
    id uuid NOT NULL,
    source_user_id uuid,
    recipient_user_id uuid,
    card_id uuid,
    amount integer DEFAULT 0 NOT NULL,
    type text NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: feature_switches; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.feature_switches (
    id uuid NOT NULL,
    name text NOT NULL,
    state text NOT NULL,
    testing_users jsonb DEFAULT '{}'::jsonb NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


--
-- Name: invites; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.invites (
    id uuid NOT NULL,
    node_id uuid NOT NULL,
    token text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    remaining_uses integer DEFAULT 0 NOT NULL,
    hide_from_user boolean DEFAULT false NOT NULL,
    group_id uuid,
    channel_id uuid,
    system_invite boolean DEFAULT false NOT NULL,
    CONSTRAINT invites_token_check CHECK ((token <> ''::text))
);


--
-- Name: leaderboard_rankings; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.leaderboard_rankings (
    user_id uuid NOT NULL,
    rank integer NOT NULL,
    coins_earned integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: mentions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.mentions (
    id uuid NOT NULL,
    in_card uuid NOT NULL,
    mentioned_user uuid NOT NULL,
    mentioned_alias uuid,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    deleted_at timestamp without time zone
);


--
-- Name: notifications; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.notifications (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    target_id uuid,
    type text NOT NULL,
    seen_at timestamp without time zone,
    opened_at timestamp without time zone,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    target_alias_id uuid
);


--
-- Name: notifications_comments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.notifications_comments (
    notification_id uuid NOT NULL,
    card_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


--
-- Name: notifications_follows; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.notifications_follows (
    notification_id uuid NOT NULL,
    follower_id uuid NOT NULL,
    followee_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL
);


--
-- Name: notifications_leaderboard_data; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.notifications_leaderboard_data (
    notification_id uuid NOT NULL,
    rank integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: notifications_mentions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.notifications_mentions (
    notification_id uuid NOT NULL,
    mention_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


--
-- Name: notifications_reactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.notifications_reactions (
    notification_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    user_id uuid NOT NULL,
    card_id uuid NOT NULL
);


--
-- Name: oauth_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.oauth_accounts (
    id uuid NOT NULL,
    provider text NOT NULL,
    subject text NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT oauth_accounts_check CHECK ((subject <> ''::text))
);


--
-- Name: popular_ranks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.popular_ranks (
    card_id uuid NOT NULL,
    views integer DEFAULT 0 NOT NULL,
    upvote_count integer DEFAULT 0 NOT NULL,
    downvote_count integer DEFAULT 0 NOT NULL,
    comment_count integer DEFAULT 0 NOT NULL,
    score_mod double precision DEFAULT 0.0 NOT NULL,
    created_at_timestamp integer DEFAULT 0 NOT NULL,
    unique_commenters_count integer DEFAULT 0 NOT NULL
);


--
-- Name: reset_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.reset_tokens (
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    expires timestamp without time zone NOT NULL,
    user_id uuid NOT NULL,
    token_hash text NOT NULL
);


--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migration (
    version character varying(255) NOT NULL
);


--
-- Name: score_modifications; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.score_modifications (
    id uuid NOT NULL,
    card_id uuid NOT NULL,
    user_id uuid NOT NULL,
    strength text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


--
-- Name: sessions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sessions (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


--
-- Name: settings; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.settings (
    id integer NOT NULL,
    signups_frozen boolean DEFAULT false NOT NULL,
    maintenance_mode boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


--
-- Name: settings_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.settings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: settings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.settings_id_seq OWNED BY public.settings.id;


--
-- Name: subscriptions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.subscriptions (
    user_id uuid NOT NULL,
    card_id uuid NOT NULL,
    type text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT subscriptions_type_check CHECK (((type = 'boost'::text) OR (type = 'comment'::text)))
);


--
-- Name: thread_mutes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.thread_mutes (
    user_id uuid NOT NULL,
    thread_root_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: user_blocks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_blocks (
    user_id uuid NOT NULL,
    blocked_user uuid,
    blocked_alias uuid,
    for_thread uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: user_card_ranks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_card_ranks (
    user_id uuid NOT NULL,
    card_id uuid NOT NULL
);


--
-- Name: user_feeds; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_feeds (
    user_id uuid NOT NULL,
    card_id uuid NOT NULL,
    "position" integer NOT NULL,
    current_top boolean DEFAULT false NOT NULL,
    rank_updated_at timestamp without time zone,
    last_visited_at timestamp without time zone,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: user_follows; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_follows (
    follower_id uuid NOT NULL,
    followee_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL
);


--
-- Name: user_mutes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_mutes (
    user_id uuid NOT NULL,
    muted_user_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: user_popular_feeds; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_popular_feeds (
    user_id uuid NOT NULL,
    card_id uuid NOT NULL,
    "position" integer NOT NULL
);


--
-- Name: user_reactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_reactions (
    user_id uuid NOT NULL,
    card_id uuid NOT NULL,
    alias_id uuid,
    type text NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT user_reactions_type_check CHECK ((type = ANY (ARRAY['like'::text, 'dislike'::text])))
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    display_name text NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    profile_image_path text NOT NULL,
    bio text NOT NULL,
    email text NOT NULL,
    password_hash text,
    password_salt text,
    username text NOT NULL,
    devices jsonb DEFAULT '{}'::jsonb NOT NULL,
    admin boolean NOT NULL,
    search_key text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    feed_updated_at timestamp without time zone,
    cover_image_path text DEFAULT ''::text NOT NULL,
    deleted_at timestamp without time zone,
    allow_email boolean DEFAULT true NOT NULL,
    joined_from_invite uuid,
    botched_signup boolean DEFAULT false,
    blocked_at timestamp without time zone,
    possible_uninstall boolean DEFAULT false NOT NULL,
    got_delayed_invites boolean DEFAULT false NOT NULL,
    shadowbanned_at timestamp without time zone,
    seen_intro_cards boolean DEFAULT false NOT NULL,
    is_default boolean DEFAULT false NOT NULL,
    feed_last_updated_at timestamp without time zone,
    disable_feed boolean DEFAULT false NOT NULL,
    coin_balance integer DEFAULT 0 NOT NULL,
    temporary_coin_balance integer DEFAULT 0 NOT NULL,
    is_verified boolean DEFAULT false NOT NULL,
    coin_reward_last_updated_at timestamp without time zone,
    is_internal boolean DEFAULT false NOT NULL,
    CONSTRAINT balance_positive CHECK ((coin_balance >= 0)),
    CONSTRAINT user_display_name_check CHECK ((display_name <> ''::text)),
    CONSTRAINT user_email_check CHECK ((email <> ''::text)),
    CONSTRAINT user_username_check CHECK ((username <> ''::text))
);


--
-- Name: waitlist; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.waitlist (
    email character varying(256) NOT NULL,
    comment text,
    created_at timestamp without time zone NOT NULL,
    name text DEFAULT ''::text NOT NULL
);


--
-- Name: settings id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.settings ALTER COLUMN id SET DEFAULT nextval('public.settings_id_seq'::regclass);


--
-- Name: activities activities_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.activities
    ADD CONSTRAINT activities_pkey PRIMARY KEY (id);


--
-- Name: announcements announcements_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.announcements
    ADD CONSTRAINT announcements_pkey PRIMARY KEY (id);


--
-- Name: anonymous_aliases anonymous_aliases_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.anonymous_aliases
    ADD CONSTRAINT anonymous_aliases_pkey PRIMARY KEY (id);


--
-- Name: cards cards_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cards
    ADD CONSTRAINT cards_pkey PRIMARY KEY (id);


--
-- Name: channel_mutes channel_mutes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_mutes
    ADD CONSTRAINT channel_mutes_pkey PRIMARY KEY (user_id, channel_id);


--
-- Name: channel_memberships channel_subscriptions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_memberships
    ADD CONSTRAINT channel_subscriptions_pkey PRIMARY KEY (user_id, channel_id);


--
-- Name: channels channels_handle_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channels
    ADD CONSTRAINT channels_handle_key UNIQUE (handle);


--
-- Name: channels channels_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channels
    ADD CONSTRAINT channels_name_key UNIQUE (name);


--
-- Name: channels channels_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channels
    ADD CONSTRAINT channels_pkey PRIMARY KEY (id);


--
-- Name: coin_rewards coin_rewards_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coin_rewards
    ADD CONSTRAINT coin_rewards_pkey PRIMARY KEY (id);


--
-- Name: coin_transactions coin_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coin_transactions
    ADD CONSTRAINT coin_transactions_pkey PRIMARY KEY (id);


--
-- Name: feature_switches feature_switches_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.feature_switches
    ADD CONSTRAINT feature_switches_pkey PRIMARY KEY (id);


--
-- Name: user_follows followers_follower_id_followee_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_follows
    ADD CONSTRAINT followers_follower_id_followee_id_key UNIQUE (follower_id, followee_id);


--
-- Name: invites invites_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.invites
    ADD CONSTRAINT invites_pkey PRIMARY KEY (id);


--
-- Name: leaderboard_rankings leaderboard_rankings_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.leaderboard_rankings
    ADD CONSTRAINT leaderboard_rankings_pkey PRIMARY KEY (user_id);


--
-- Name: mentions mentions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.mentions
    ADD CONSTRAINT mentions_pkey PRIMARY KEY (id);


--
-- Name: notifications_comments notifications_comments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications_comments
    ADD CONSTRAINT notifications_comments_pkey PRIMARY KEY (notification_id, card_id);


--
-- Name: notifications_follows notifications_follows_follower_id_followee_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications_follows
    ADD CONSTRAINT notifications_follows_follower_id_followee_id_key UNIQUE (follower_id, followee_id);


--
-- Name: notifications_leaderboard_data notifications_leaderboard_data_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications_leaderboard_data
    ADD CONSTRAINT notifications_leaderboard_data_pkey PRIMARY KEY (notification_id);


--
-- Name: notifications_mentions notifications_mentions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications_mentions
    ADD CONSTRAINT notifications_mentions_pkey PRIMARY KEY (notification_id, mention_id);


--
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--
-- Name: notifications_reactions notifications_reactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications_reactions
    ADD CONSTRAINT notifications_reactions_pkey PRIMARY KEY (notification_id, user_id, card_id);


--
-- Name: oauth_accounts oauth_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.oauth_accounts
    ADD CONSTRAINT oauth_accounts_pkey PRIMARY KEY (id);


--
-- Name: oauth_accounts oauth_accounts_provider_subject_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.oauth_accounts
    ADD CONSTRAINT oauth_accounts_provider_subject_key UNIQUE (provider, subject);


--
-- Name: popular_ranks popular_ranks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.popular_ranks
    ADD CONSTRAINT popular_ranks_pkey PRIMARY KEY (card_id);


--
-- Name: reset_tokens reset_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reset_tokens
    ADD CONSTRAINT reset_tokens_pkey PRIMARY KEY (user_id);


--
-- Name: score_modifications score_modifications_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.score_modifications
    ADD CONSTRAINT score_modifications_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: settings settings_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.settings
    ADD CONSTRAINT settings_pkey PRIMARY KEY (id);


--
-- Name: subscriptions subscriptions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_pkey PRIMARY KEY (user_id, card_id, type);


--
-- Name: thread_mutes thread_mutes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.thread_mutes
    ADD CONSTRAINT thread_mutes_pkey PRIMARY KEY (user_id, thread_root_id);


--
-- Name: user_blocks user_blocks_user_id_blocked_alias_for_thread_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_blocks
    ADD CONSTRAINT user_blocks_user_id_blocked_alias_for_thread_key UNIQUE (user_id, blocked_alias, for_thread);


--
-- Name: user_blocks user_blocks_user_id_blocked_user_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_blocks
    ADD CONSTRAINT user_blocks_user_id_blocked_user_key UNIQUE (user_id, blocked_user);


--
-- Name: user_card_ranks user_card_ranks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_card_ranks
    ADD CONSTRAINT user_card_ranks_pkey PRIMARY KEY (user_id, card_id);


--
-- Name: user_feeds user_feeds_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_feeds
    ADD CONSTRAINT user_feeds_pkey PRIMARY KEY (user_id, card_id);


--
-- Name: user_feeds user_feeds_user_id_position_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_feeds
    ADD CONSTRAINT user_feeds_user_id_position_key UNIQUE (user_id, "position");


--
-- Name: user_mutes user_mutes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_mutes
    ADD CONSTRAINT user_mutes_pkey PRIMARY KEY (user_id, muted_user_id);


--
-- Name: user_popular_feeds user_popular_feeds_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_popular_feeds
    ADD CONSTRAINT user_popular_feeds_pkey PRIMARY KEY (user_id, card_id);


--
-- Name: user_reactions user_reactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_reactions
    ADD CONSTRAINT user_reactions_pkey PRIMARY KEY (user_id, card_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: waitlist waitlist_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.waitlist
    ADD CONSTRAINT waitlist_pkey PRIMARY KEY (email);


--
-- Name: announcements_deleted_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX announcements_deleted_at_idx ON public.announcements USING btree (deleted_at);


--
-- Name: anonymous_aliases_username_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX anonymous_aliases_username_idx ON public.anonymous_aliases USING btree (username);


--
-- Name: cards_created_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX cards_created_at_idx ON public.cards USING btree (created_at);


--
-- Name: cards_deleted_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX cards_deleted_at_idx ON public.cards USING btree (deleted_at);


--
-- Name: cards_thread_reply_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX cards_thread_reply_id_idx ON public.cards USING btree (thread_reply_id);


--
-- Name: invites_token_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX invites_token_idx ON public.invites USING btree (token);


--
-- Name: oauth_accounts_subject_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX oauth_accounts_subject_idx ON public.oauth_accounts USING btree (subject);


--
-- Name: users_deleted_at_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX users_deleted_at_idx ON public.users USING btree (deleted_at);


--
-- Name: users_email_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX users_email_idx ON public.users USING btree (email);


--
-- Name: users_username_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX users_username_idx ON public.users USING btree (username);


--
-- Name: version_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX version_idx ON public.schema_migration USING btree (version);


--
-- Name: activities activities_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.activities
    ADD CONSTRAINT activities_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: cards cards_alias_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cards
    ADD CONSTRAINT cards_alias_id_fkey FOREIGN KEY (alias_id) REFERENCES public.anonymous_aliases(id);


--
-- Name: cards cards_channel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cards
    ADD CONSTRAINT cards_channel_id_fkey FOREIGN KEY (channel_id) REFERENCES public.channels(id);


--
-- Name: cards cards_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cards
    ADD CONSTRAINT cards_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public.users(id);


--
-- Name: cards cards_thread_reply_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cards
    ADD CONSTRAINT cards_thread_reply_id_fkey FOREIGN KEY (thread_reply_id) REFERENCES public.cards(id);


--
-- Name: cards cards_thread_root_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cards
    ADD CONSTRAINT cards_thread_root_id_fkey FOREIGN KEY (thread_root_id) REFERENCES public.cards(id);


--
-- Name: channel_mutes channel_mutes_channel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_mutes
    ADD CONSTRAINT channel_mutes_channel_id_fkey FOREIGN KEY (channel_id) REFERENCES public.channels(id);


--
-- Name: channel_mutes channel_mutes_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_mutes
    ADD CONSTRAINT channel_mutes_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: channel_memberships channel_subscriptions_channel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_memberships
    ADD CONSTRAINT channel_subscriptions_channel_id_fkey FOREIGN KEY (channel_id) REFERENCES public.channels(id);


--
-- Name: channel_memberships channel_subscriptions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channel_memberships
    ADD CONSTRAINT channel_subscriptions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: channels channels_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.channels
    ADD CONSTRAINT channels_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public.users(id);


--
-- Name: coin_rewards coin_rewards_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coin_rewards
    ADD CONSTRAINT coin_rewards_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: coin_transactions coin_transactions_card_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coin_transactions
    ADD CONSTRAINT coin_transactions_card_id_fkey FOREIGN KEY (card_id) REFERENCES public.cards(id);


--
-- Name: coin_transactions coin_transactions_recipient_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coin_transactions
    ADD CONSTRAINT coin_transactions_recipient_user_id_fkey FOREIGN KEY (recipient_user_id) REFERENCES public.users(id);


--
-- Name: coin_transactions coin_transactions_source_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coin_transactions
    ADD CONSTRAINT coin_transactions_source_user_id_fkey FOREIGN KEY (source_user_id) REFERENCES public.users(id);


--
-- Name: user_follows followers_followee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_follows
    ADD CONSTRAINT followers_followee_id_fkey FOREIGN KEY (followee_id) REFERENCES public.users(id);


--
-- Name: user_follows followers_follower_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_follows
    ADD CONSTRAINT followers_follower_id_fkey FOREIGN KEY (follower_id) REFERENCES public.users(id);


--
-- Name: invites invites_channel_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.invites
    ADD CONSTRAINT invites_channel_id_fkey FOREIGN KEY (channel_id) REFERENCES public.channels(id);


--
-- Name: invites invites_node_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.invites
    ADD CONSTRAINT invites_node_id_fkey FOREIGN KEY (node_id) REFERENCES public.users(id);


--
-- Name: leaderboard_rankings leaderboard_rankings_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.leaderboard_rankings
    ADD CONSTRAINT leaderboard_rankings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: mentions mentions_in_card_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.mentions
    ADD CONSTRAINT mentions_in_card_fkey FOREIGN KEY (in_card) REFERENCES public.cards(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: mentions mentions_mentioned_alias_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.mentions
    ADD CONSTRAINT mentions_mentioned_alias_fkey FOREIGN KEY (mentioned_alias) REFERENCES public.anonymous_aliases(id) ON UPDATE CASCADE ON DELETE SET NULL;


--
-- Name: mentions mentions_mentioned_user_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.mentions
    ADD CONSTRAINT mentions_mentioned_user_fkey FOREIGN KEY (mentioned_user) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: notifications_comments notifications_comments_card_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications_comments
    ADD CONSTRAINT notifications_comments_card_id_fkey FOREIGN KEY (card_id) REFERENCES public.cards(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: notifications_comments notifications_comments_notification_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications_comments
    ADD CONSTRAINT notifications_comments_notification_id_fkey FOREIGN KEY (notification_id) REFERENCES public.notifications(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: notifications_follows notifications_follows_followee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications_follows
    ADD CONSTRAINT notifications_follows_followee_id_fkey FOREIGN KEY (followee_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: notifications_follows notifications_follows_follower_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications_follows
    ADD CONSTRAINT notifications_follows_follower_id_fkey FOREIGN KEY (follower_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: notifications_follows notifications_follows_follower_id_fkey1; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications_follows
    ADD CONSTRAINT notifications_follows_follower_id_fkey1 FOREIGN KEY (follower_id, followee_id) REFERENCES public.user_follows(follower_id, followee_id) ON DELETE CASCADE;


--
-- Name: notifications_follows notifications_follows_notification_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications_follows
    ADD CONSTRAINT notifications_follows_notification_id_fkey FOREIGN KEY (notification_id) REFERENCES public.notifications(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: notifications_leaderboard_data notifications_leaderboard_data_notification_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications_leaderboard_data
    ADD CONSTRAINT notifications_leaderboard_data_notification_id_fkey FOREIGN KEY (notification_id) REFERENCES public.notifications(id);


--
-- Name: notifications_mentions notifications_mentions_mention_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications_mentions
    ADD CONSTRAINT notifications_mentions_mention_id_fkey FOREIGN KEY (mention_id) REFERENCES public.mentions(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: notifications_mentions notifications_mentions_notification_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications_mentions
    ADD CONSTRAINT notifications_mentions_notification_id_fkey FOREIGN KEY (notification_id) REFERENCES public.notifications(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: notifications_reactions notifications_reactions_notification_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications_reactions
    ADD CONSTRAINT notifications_reactions_notification_id_fkey FOREIGN KEY (notification_id) REFERENCES public.notifications(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: notifications_reactions notifications_reactions_user_reactions_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications_reactions
    ADD CONSTRAINT notifications_reactions_user_reactions_fkey FOREIGN KEY (user_id, card_id) REFERENCES public.user_reactions(user_id, card_id) ON DELETE CASCADE;


--
-- Name: notifications notifications_target_alias_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_target_alias_id_fk FOREIGN KEY (target_alias_id) REFERENCES public.anonymous_aliases(id);


--
-- Name: oauth_accounts oauth_accounts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.oauth_accounts
    ADD CONSTRAINT oauth_accounts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: popular_ranks popular_ranks_card_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.popular_ranks
    ADD CONSTRAINT popular_ranks_card_id_fkey FOREIGN KEY (card_id) REFERENCES public.cards(id);


--
-- Name: reset_tokens reset_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reset_tokens
    ADD CONSTRAINT reset_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: sessions sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: subscriptions subscriptions_card_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_card_id_fkey FOREIGN KEY (card_id) REFERENCES public.cards(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: subscriptions subscriptions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: thread_mutes thread_mutes_thread_root_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.thread_mutes
    ADD CONSTRAINT thread_mutes_thread_root_id_fkey FOREIGN KEY (thread_root_id) REFERENCES public.cards(id);


--
-- Name: thread_mutes thread_mutes_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.thread_mutes
    ADD CONSTRAINT thread_mutes_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_blocks user_blocks_blocked_alias_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_blocks
    ADD CONSTRAINT user_blocks_blocked_alias_fkey FOREIGN KEY (blocked_alias) REFERENCES public.anonymous_aliases(id);


--
-- Name: user_blocks user_blocks_blocked_user_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_blocks
    ADD CONSTRAINT user_blocks_blocked_user_fkey FOREIGN KEY (blocked_user) REFERENCES public.users(id);


--
-- Name: user_blocks user_blocks_for_thread_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_blocks
    ADD CONSTRAINT user_blocks_for_thread_fkey FOREIGN KEY (for_thread) REFERENCES public.cards(id);


--
-- Name: user_blocks user_blocks_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_blocks
    ADD CONSTRAINT user_blocks_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_card_ranks user_card_ranks_card_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_card_ranks
    ADD CONSTRAINT user_card_ranks_card_id_fkey FOREIGN KEY (card_id) REFERENCES public.cards(id);


--
-- Name: user_card_ranks user_card_ranks_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_card_ranks
    ADD CONSTRAINT user_card_ranks_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_feeds user_feeds_card_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_feeds
    ADD CONSTRAINT user_feeds_card_id_fkey FOREIGN KEY (card_id) REFERENCES public.cards(id);


--
-- Name: user_feeds user_feeds_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_feeds
    ADD CONSTRAINT user_feeds_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_mutes user_mutes_muted_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_mutes
    ADD CONSTRAINT user_mutes_muted_user_id_fkey FOREIGN KEY (muted_user_id) REFERENCES public.users(id);


--
-- Name: user_mutes user_mutes_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_mutes
    ADD CONSTRAINT user_mutes_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_popular_feeds user_popular_feeds_card_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_popular_feeds
    ADD CONSTRAINT user_popular_feeds_card_id_fkey FOREIGN KEY (card_id) REFERENCES public.cards(id);


--
-- Name: user_popular_feeds user_popular_feeds_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_popular_feeds
    ADD CONSTRAINT user_popular_feeds_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: user_reactions user_reactions_alias_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_reactions
    ADD CONSTRAINT user_reactions_alias_id_fkey FOREIGN KEY (alias_id) REFERENCES public.anonymous_aliases(id);


--
-- Name: user_reactions user_reactions_card_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_reactions
    ADD CONSTRAINT user_reactions_card_id_fkey FOREIGN KEY (card_id) REFERENCES public.cards(id);


--
-- Name: user_reactions user_reactions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_reactions
    ADD CONSTRAINT user_reactions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: users users_invite_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_invite_id_fkey FOREIGN KEY (joined_from_invite) REFERENCES public.invites(id);


--
-- PostgreSQL database dump complete
--

