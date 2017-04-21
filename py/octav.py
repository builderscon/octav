"""OCTAV Client Library"""
"""DO NOT EDIT: This file was generated from ../spec/v2/api.json on Fri Apr 21 07:10:38 2017"""

import certifi
import feedparser
import functools
import json
import os
import re
import time
import urllib3

if os.getenv('SERVER_SOFTWARE', '').startswith('Google App Engine/') or os.getenv('SERVER_SOFTWARE', '').startswith('Development/'):
    from urllib3.contrib.appengine import AppEngineManager as PoolManager
else:
    from urllib3 import PoolManager

import sys
if sys.version[0] == "3":
    from urllib.parse import urlencode
else:
    from urllib import urlencode

class MissingRequiredArgument(Exception):
    pass

def parse_rfc3339(s):
    return time.mktime(feedparser._parse_date(s))


class Session(object):
  def __init__(self, client, f, sid, expires):
    if not f:
      raise MissingRequiredArgument('f is required')
    if not client:
      raise MissingRequiredArgument('client is required')
    if not sid:
      raise MissingRequiredArgument('sid is required')
    if not expires:
      raise MissingRequiredArgument('expires is required')

    self.update_func = f
    self.sid = sid
    self.expires = expires
    self.client = client

  def last_error(self):
    return self.client.last_error()

  # renews the octav session. returns false if there was no need to
  # renew, true if the session was renewed. None is returned if
  # there was an error
  def renew(self):
    if self.expires > time.mktime(time.gmtime()):
        return False
    s = self.update_func()
    if s is None:
        return None
    self.sid = s.get('sid')
    self.expires = parse_rfc3339(s.get('expires'))
    return True

  def add_conference_admin (self, admin_id, conference_id, extra_headers=None):
    self.renew()
    return self.client.add_conference_admin(admin_id, conference_id, extra_headers={'X-Octav-Session-ID': self.sid})

  def add_conference_credential (self, conference_id, data, type, extra_headers=None):
    self.renew()
    return self.client.add_conference_credential(conference_id, data, type, extra_headers={'X-Octav-Session-ID': self.sid})

  def add_conference_date (self, conference_id, date, extra_headers=None):
    self.renew()
    return self.client.add_conference_date(conference_id, date, extra_headers={'X-Octav-Session-ID': self.sid})

  def add_conference_series_admin (self, admin_id, series_id, extra_headers=None):
    self.renew()
    return self.client.add_conference_series_admin(admin_id, series_id, extra_headers={'X-Octav-Session-ID': self.sid})

  def add_conference_staff (self, conference_id, staff_id, extra_headers=None):
    self.renew()
    return self.client.add_conference_staff(conference_id, staff_id, extra_headers={'X-Octav-Session-ID': self.sid})

  def add_conference_venue (self, conference_id, venue_id, extra_headers=None):
    self.renew()
    return self.client.add_conference_venue(conference_id, venue_id, extra_headers={'X-Octav-Session-ID': self.sid})

  def add_featured_speaker (self, conference_id, description, display_name, avatar_url=None, speaker_id=None, extra_headers=None, **args):
    self.renew()
    return self.client.add_featured_speaker(conference_id, description, display_name, avatar_url=avatar_url, speaker_id=speaker_id, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def add_session_type (self, abstract, conference_id, duration, name, submission_end=None, submission_start=None, extra_headers=None, **args):
    self.renew()
    return self.client.add_session_type(abstract, conference_id, duration, name, submission_end=submission_end, submission_start=submission_start, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def add_sponsor (self, conference_id, group_name, name, url, logo_url=None, sort_order=None, extra_headers=None, **args):
    self.renew()
    return self.client.add_sponsor(conference_id, group_name, name, url, logo_url=logo_url, sort_order=sort_order, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def confirm_temporary_email (self, confirmation_key, target_id, extra_headers=None):
    self.renew()
    return self.client.confirm_temporary_email(confirmation_key, target_id, extra_headers={'X-Octav-Session-ID': self.sid})

  def create_blog_entry (self, conference_id, url, status=None, title=None, extra_headers=None):
    self.renew()
    return self.client.create_blog_entry(conference_id, url, status=status, title=title, extra_headers={'X-Octav-Session-ID': self.sid})

  def create_client_session (self, access_token, auth_via, extra_headers=None):
    self.renew()
    return self.client.create_client_session(access_token, auth_via, extra_headers={'X-Octav-Session-ID': self.sid})

  def create_conference (self, series_id, slug, title, cfp_lead_text=None, cfp_post_submit_instructions=None, cfp_pre_submit_instructions=None, contact_information=None, description=None, sub_title=None, timezone=None, extra_headers=None, **args):
    self.renew()
    return self.client.create_conference(series_id, slug, title, cfp_lead_text=cfp_lead_text, cfp_post_submit_instructions=cfp_post_submit_instructions, cfp_pre_submit_instructions=cfp_pre_submit_instructions, contact_information=contact_information, description=description, sub_title=sub_title, timezone=timezone, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def create_conference_series (self, session, slug, title, description=None, sid=None, extra_headers=None):
    self.renew()
    return self.client.create_conference_series(session, slug, title, description=description, sid=sid, extra_headers={'X-Octav-Session-ID': self.sid})

  def create_external_resource (self, conference_id, title, url, description=None, image_url=None, sort_order=None, extra_headers=None, **args):
    self.renew()
    return self.client.create_external_resource(conference_id, title, url, description=description, image_url=image_url, sort_order=sort_order, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def create_question (self, body, session_id, extra_headers=None):
    self.renew()
    return self.client.create_question(body, session_id, extra_headers={'X-Octav-Session-ID': self.sid})

  def create_room (self, name, venue_id, capacity=None, extra_headers=None, **args):
    self.renew()
    return self.client.create_room(name, venue_id, capacity=capacity, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def create_session (self, conference_id, session_type_id, speaker_id, abstract=None, category=None, material_level=None, materials_release=None, memo=None, photo_release=None, recording_release=None, slide_language=None, slide_subtitles=None, slide_url=None, spoken_language=None, tags=None, title=None, video_url=None, extra_headers=None, **args):
    self.renew()
    return self.client.create_session(conference_id, session_type_id, speaker_id, abstract=abstract, category=category, material_level=material_level, materials_release=materials_release, memo=memo, photo_release=photo_release, recording_release=recording_release, slide_language=slide_language, slide_subtitles=slide_subtitles, slide_url=slide_url, spoken_language=spoken_language, tags=tags, title=title, video_url=video_url, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def create_session_survey_response (self, material_quality, overall_rating, session_id, speaker_knowledge, speaker_presentation, user_prior_knowledge, comment_good=None, comment_improvement=None, extra_headers=None):
    self.renew()
    return self.client.create_session_survey_response(material_quality, overall_rating, session_id, speaker_knowledge, speaker_presentation, user_prior_knowledge, comment_good=comment_good, comment_improvement=comment_improvement, extra_headers={'X-Octav-Session-ID': self.sid})

  def create_temporary_email (self, email, target_id, extra_headers=None):
    self.renew()
    return self.client.create_temporary_email(email, target_id, extra_headers={'X-Octav-Session-ID': self.sid})

  def create_track (self, conference_id, room_id, name=None, sort_order=None, extra_headers=None, **args):
    self.renew()
    return self.client.create_track(conference_id, room_id, name=name, sort_order=sort_order, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def create_user (self, auth_user_id, auth_via, nickname, avatar_url=None, email=None, first_name=None, lang=None, last_name=None, tshirt_size=None, extra_headers=None, **args):
    self.renew()
    return self.client.create_user(auth_user_id, auth_via, nickname, avatar_url=avatar_url, email=email, first_name=first_name, lang=lang, last_name=last_name, tshirt_size=tshirt_size, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def create_venue (self, address, name, latitude=None, longitude=None, place_id=None, url=None, extra_headers=None, **args):
    self.renew()
    return self.client.create_venue(address, name, latitude=latitude, longitude=longitude, place_id=place_id, url=url, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def delete_blog_entry (self, id, extra_headers=None):
    self.renew()
    return self.client.delete_blog_entry(id, extra_headers={'X-Octav-Session-ID': self.sid})

  def delete_conference (self, id, extra_headers=None):
    self.renew()
    return self.client.delete_conference(id, extra_headers={'X-Octav-Session-ID': self.sid})

  def delete_conference_admin (self, admin_id, conference_id, extra_headers=None):
    self.renew()
    return self.client.delete_conference_admin(admin_id, conference_id, extra_headers={'X-Octav-Session-ID': self.sid})

  def delete_conference_date (self, conference_id, date, extra_headers=None):
    self.renew()
    return self.client.delete_conference_date(conference_id, date, extra_headers={'X-Octav-Session-ID': self.sid})

  def delete_conference_series (self, id, extra_headers=None):
    self.renew()
    return self.client.delete_conference_series(id, extra_headers={'X-Octav-Session-ID': self.sid})

  def delete_conference_staff (self, conference_id, staff_id, extra_headers=None):
    self.renew()
    return self.client.delete_conference_staff(conference_id, staff_id, extra_headers={'X-Octav-Session-ID': self.sid})

  def delete_conference_venue (self, conference_id, venue_id, extra_headers=None):
    self.renew()
    return self.client.delete_conference_venue(conference_id, venue_id, extra_headers={'X-Octav-Session-ID': self.sid})

  def delete_external_resource (self, id, extra_headers=None):
    self.renew()
    return self.client.delete_external_resource(id, extra_headers={'X-Octav-Session-ID': self.sid})

  def delete_featured_speaker (self, id, extra_headers=None):
    self.renew()
    return self.client.delete_featured_speaker(id, extra_headers={'X-Octav-Session-ID': self.sid})

  def delete_question (self, id, extra_headers=None):
    self.renew()
    return self.client.delete_question(id, extra_headers={'X-Octav-Session-ID': self.sid})

  def delete_room (self, id, extra_headers=None):
    self.renew()
    return self.client.delete_room(id, extra_headers={'X-Octav-Session-ID': self.sid})

  def delete_session (self, id, extra_headers=None):
    self.renew()
    return self.client.delete_session(id, extra_headers={'X-Octav-Session-ID': self.sid})

  def delete_session_type (self, id, extra_headers=None):
    self.renew()
    return self.client.delete_session_type(id, extra_headers={'X-Octav-Session-ID': self.sid})

  def delete_sponsor (self, id, extra_headers=None):
    self.renew()
    return self.client.delete_sponsor(id, extra_headers={'X-Octav-Session-ID': self.sid})

  def delete_track (self, id, extra_headers=None):
    self.renew()
    return self.client.delete_track(id, extra_headers={'X-Octav-Session-ID': self.sid})

  def delete_user (self, id, extra_headers=None):
    self.renew()
    return self.client.delete_user(id, extra_headers={'X-Octav-Session-ID': self.sid})

  def delete_venue (self, id, extra_headers=None):
    self.renew()
    return self.client.delete_venue(id, extra_headers={'X-Octav-Session-ID': self.sid})

  def get_conference_schedule (self, conference_id, lang=None, extra_headers=None):
    self.renew()
    return self.client.get_conference_schedule(conference_id, lang=lang, extra_headers={'X-Octav-Session-ID': self.sid})

  def health_check (self, extra_headers=None):
    self.renew()
    return self.client.health_check(extra_headers={'X-Octav-Session-ID': self.sid})

  def list_blog_entries (self, conference_id=None, lang=None, status=None, extra_headers=None):
    self.renew()
    return self.client.list_blog_entries(conference_id=conference_id, lang=lang, status=status, extra_headers={'X-Octav-Session-ID': self.sid})

  def list_conference (self, lang=None, limit=None, organizers=None, range_end=None, range_start=None, since=None, status=None, extra_headers=None):
    self.renew()
    return self.client.list_conference(lang=lang, limit=limit, organizers=organizers, range_end=range_end, range_start=range_start, since=since, status=status, extra_headers={'X-Octav-Session-ID': self.sid})

  def list_conference_admin (self, conference_id, extra_headers=None):
    self.renew()
    return self.client.list_conference_admin(conference_id, extra_headers={'X-Octav-Session-ID': self.sid})

  def list_conference_date (self, conference_id, extra_headers=None):
    self.renew()
    return self.client.list_conference_date(conference_id, extra_headers={'X-Octav-Session-ID': self.sid})

  def list_conference_series (self, limit=None, since=None, extra_headers=None):
    self.renew()
    return self.client.list_conference_series(limit=limit, since=since, extra_headers={'X-Octav-Session-ID': self.sid})

  def list_conference_staff (self, conference_id=None, lang=None, extra_headers=None):
    self.renew()
    return self.client.list_conference_staff(conference_id=conference_id, lang=lang, extra_headers={'X-Octav-Session-ID': self.sid})

  def list_conferences_by_organizer (self, lang=None, limit=None, organizer_id=None, since=None, status=None, extra_headers=None):
    self.renew()
    return self.client.list_conferences_by_organizer(lang=lang, limit=limit, organizer_id=organizer_id, since=since, status=status, extra_headers={'X-Octav-Session-ID': self.sid})

  def list_external_resource (self, conference_id, lang=None, limit=None, since=None, extra_headers=None):
    self.renew()
    return self.client.list_external_resource(conference_id, lang=lang, limit=limit, since=since, extra_headers={'X-Octav-Session-ID': self.sid})

  def list_featured_speakers (self, conference_id=None, lang=None, limit=None, since=None, extra_headers=None):
    self.renew()
    return self.client.list_featured_speakers(conference_id=conference_id, lang=lang, limit=limit, since=since, extra_headers={'X-Octav-Session-ID': self.sid})

  def list_question (self, session_id, limit=None, since=None, extra_headers=None):
    self.renew()
    return self.client.list_question(session_id, limit=limit, since=since, extra_headers={'X-Octav-Session-ID': self.sid})

  def list_room (self, venue_id, lang=None, limit=None, extra_headers=None):
    self.renew()
    return self.client.list_room(venue_id, lang=lang, limit=limit, extra_headers={'X-Octav-Session-ID': self.sid})

  def list_session_types_by_conference (self, conference_id=None, lang=None, extra_headers=None):
    self.renew()
    return self.client.list_session_types_by_conference(conference_id=conference_id, lang=lang, extra_headers={'X-Octav-Session-ID': self.sid})

  def list_sessions (self, conference_id=None, confirmed=None, lang=None, limit=None, range_end=None, range_start=None, since=None, speaker_id=None, status=None, extra_headers=None):
    self.renew()
    return self.client.list_sessions(conference_id=conference_id, confirmed=confirmed, lang=lang, limit=limit, range_end=range_end, range_start=range_start, since=since, speaker_id=speaker_id, status=status, extra_headers={'X-Octav-Session-ID': self.sid})

  def list_sponsors (self, conference_id=None, lang=None, limit=None, since=None, extra_headers=None):
    self.renew()
    return self.client.list_sponsors(conference_id=conference_id, lang=lang, limit=limit, since=since, extra_headers={'X-Octav-Session-ID': self.sid})

  def list_user (self, lang=None, limit=None, pattern=None, since=None, extra_headers=None):
    self.renew()
    return self.client.list_user(lang=lang, limit=limit, pattern=pattern, since=since, extra_headers={'X-Octav-Session-ID': self.sid})

  def list_venue (self, lang=None, limit=None, since=None, extra_headers=None):
    self.renew()
    return self.client.list_venue(lang=lang, limit=limit, since=since, extra_headers={'X-Octav-Session-ID': self.sid})

  def lookup_blog_entry (self, id, lang=None, extra_headers=None):
    self.renew()
    return self.client.lookup_blog_entry(id, lang=lang, extra_headers={'X-Octav-Session-ID': self.sid})

  def lookup_conference (self, id, lang=None, extra_headers=None):
    self.renew()
    return self.client.lookup_conference(id, lang=lang, extra_headers={'X-Octav-Session-ID': self.sid})

  def lookup_conference_by_slug (self, slug, lang=None, extra_headers=None):
    self.renew()
    return self.client.lookup_conference_by_slug(slug, lang=lang, extra_headers={'X-Octav-Session-ID': self.sid})

  def lookup_conference_series (self, id, lang=None, extra_headers=None):
    self.renew()
    return self.client.lookup_conference_series(id, lang=lang, extra_headers={'X-Octav-Session-ID': self.sid})

  def lookup_external_resource (self, id, lang=None, limit=None, since=None, extra_headers=None):
    self.renew()
    return self.client.lookup_external_resource(id, lang=lang, limit=limit, since=since, extra_headers={'X-Octav-Session-ID': self.sid})

  def lookup_featured_speaker (self, id, lang=None, extra_headers=None):
    self.renew()
    return self.client.lookup_featured_speaker(id, lang=lang, extra_headers={'X-Octav-Session-ID': self.sid})

  def lookup_room (self, id, extra_headers=None):
    self.renew()
    return self.client.lookup_room(id, extra_headers={'X-Octav-Session-ID': self.sid})

  def lookup_session (self, id, lang=None, extra_headers=None):
    self.renew()
    return self.client.lookup_session(id, lang=lang, extra_headers={'X-Octav-Session-ID': self.sid})

  def lookup_session_type (self, id, lang=None, extra_headers=None):
    self.renew()
    return self.client.lookup_session_type(id, lang=lang, extra_headers={'X-Octav-Session-ID': self.sid})

  def lookup_sponsor (self, id, lang=None, extra_headers=None):
    self.renew()
    return self.client.lookup_sponsor(id, lang=lang, extra_headers={'X-Octav-Session-ID': self.sid})

  def lookup_track (self, id, lang=None, extra_headers=None):
    self.renew()
    return self.client.lookup_track(id, lang=lang, extra_headers={'X-Octav-Session-ID': self.sid})

  def lookup_user (self, id, sid=None, extra_headers=None):
    self.renew()
    return self.client.lookup_user(id, sid=sid, extra_headers={'X-Octav-Session-ID': self.sid})

  def lookup_user_by_auth_user_id (self, auth_user_id, auth_via, extra_headers=None):
    self.renew()
    return self.client.lookup_user_by_auth_user_id(auth_user_id, auth_via, extra_headers={'X-Octav-Session-ID': self.sid})

  def lookup_venue (self, id, lang=None, extra_headers=None):
    self.renew()
    return self.client.lookup_venue(id, lang=lang, extra_headers={'X-Octav-Session-ID': self.sid})

  def send_all_selection_result_notification (self, conference_id, force=None, extra_headers=None):
    self.renew()
    return self.client.send_all_selection_result_notification(conference_id, force=force, extra_headers={'X-Octav-Session-ID': self.sid})

  def send_selection_result_notification (self, id, force=None, session_id=None, extra_headers=None):
    self.renew()
    return self.client.send_selection_result_notification(id, force=force, session_id=session_id, extra_headers={'X-Octav-Session-ID': self.sid})

  def set_session_video_cover (self, id, extra_headers=None):
    self.renew()
    return self.client.set_session_video_cover(id, extra_headers={'X-Octav-Session-ID': self.sid})

  def tweet_as_conference (self, conference_id, tweet, extra_headers=None):
    self.renew()
    return self.client.tweet_as_conference(conference_id, tweet, extra_headers={'X-Octav-Session-ID': self.sid})

  def update_blog_entry (self, id, status=None, title=None, url=None, extra_headers=None):
    self.renew()
    return self.client.update_blog_entry(id, status=status, title=title, url=url, extra_headers={'X-Octav-Session-ID': self.sid})

  def update_conference (self, id, cfp_lead_text=None, cfp_post_submit_instructions=None, cfp_pre_submit_instructions=None, contact_information=None, description=None, redirect_url=None, slug=None, status=None, sub_title=None, timetable_available=None, timezone=None, title=None, extra_headers=None, **args):
    self.renew()
    return self.client.update_conference(id, cfp_lead_text=cfp_lead_text, cfp_post_submit_instructions=cfp_post_submit_instructions, cfp_pre_submit_instructions=cfp_pre_submit_instructions, contact_information=contact_information, description=description, redirect_url=redirect_url, slug=slug, status=status, sub_title=sub_title, timetable_available=timetable_available, timezone=timezone, title=title, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def update_external_resource (self, id, description=None, image_url=None, sort_order=None, title=None, url=None, extra_headers=None, **args):
    self.renew()
    return self.client.update_external_resource(id, description=description, image_url=image_url, sort_order=sort_order, title=title, url=url, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def update_featured_speaker (self, id, avatar_url=None, description=None, display_name=None, speaker_id=None, extra_headers=None, **args):
    self.renew()
    return self.client.update_featured_speaker(id, avatar_url=avatar_url, description=description, display_name=display_name, speaker_id=speaker_id, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def update_room (self, id, capacity=None, name=None, venue_id=None, extra_headers=None, **args):
    self.renew()
    return self.client.update_room(id, capacity=capacity, name=name, venue_id=venue_id, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def update_session (self, id, abstract=None, category=None, conference_id=None, confirmed=None, duration=None, has_interpretation=None, material_level=None, materials_release=None, memo=None, photo_release=None, recording_release=None, session_type_id=None, slide_language=None, slide_subtitles=None, slide_url=None, sort_order=None, speaker_id=None, spoken_language=None, starts_on=None, status=None, tags=None, title=None, video_url=None, extra_headers=None, **args):
    self.renew()
    return self.client.update_session(id, abstract=abstract, category=category, conference_id=conference_id, confirmed=confirmed, duration=duration, has_interpretation=has_interpretation, material_level=material_level, materials_release=materials_release, memo=memo, photo_release=photo_release, recording_release=recording_release, session_type_id=session_type_id, slide_language=slide_language, slide_subtitles=slide_subtitles, slide_url=slide_url, sort_order=sort_order, speaker_id=speaker_id, spoken_language=spoken_language, starts_on=starts_on, status=status, tags=tags, title=title, video_url=video_url, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def update_session_type (self, id, abstract=None, duration=None, is_default=None, name=None, submission_end=None, submission_start=None, extra_headers=None, **args):
    self.renew()
    return self.client.update_session_type(id, abstract=abstract, duration=duration, is_default=is_default, name=name, submission_end=submission_end, submission_start=submission_start, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def update_sponsor (self, id, group_name=None, logo_url=None, name=None, sort_order=None, url=None, extra_headers=None, **args):
    self.renew()
    return self.client.update_sponsor(id, group_name=group_name, logo_url=logo_url, name=name, sort_order=sort_order, url=url, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def update_track (self, id, name=None, room_id=None, sort_order=None, extra_headers=None, **args):
    self.renew()
    return self.client.update_track(id, name=name, room_id=room_id, sort_order=sort_order, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def update_user (self, id, email=None, first_name=None, lang=None, last_name=None, nickname=None, tshirt_size=None, extra_headers=None, **args):
    self.renew()
    return self.client.update_user(id, email=email, first_name=first_name, lang=lang, last_name=last_name, nickname=nickname, tshirt_size=tshirt_size, extra_headers={'X-Octav-Session-ID': self.sid}, **args)

  def update_venue (self, id, latitude=None, longitude=None, name=None, place_id=None, url=None, extra_headers=None):
    self.renew()
    return self.client.update_venue(id, latitude=latitude, longitude=longitude, name=name, place_id=place_id, url=url, extra_headers={'X-Octav-Session-ID': self.sid})

  def verify_user (self, id, extra_headers=None):
    self.renew()
    return self.client.verify_user(id, extra_headers={'X-Octav-Session-ID': self.sid})


class Octav(object):
  def __init__(self, endpoint, key, secret, debug=False):
    if not endpoint:
      raise MissingRequiredArgument('endpoint is required')
    if not key:
      raise MissingRequiredArgument('key is required')
    if not secret:
      raise MissingRequiredArgument('secret is required')
    self.debug = debug
    self.endpoint = endpoint
    self.error = None
    self.http = PoolManager(cert_reqs='CERT_REQUIRED', ca_certs=certifi.where())
    self.key = key
    self.secret = secret

  def new_session(self, access_token, auth_via):
    if not access_token:
      raise MissingRequiredArgument('access_token is required')
    if not auth_via:
      raise MissingRequiredArgument('auth_via is required')

    f = functools.partial(self.create_client_session, access_token, auth_via)
    s = f()
    if s is None:
        return None
    return Session(self, f, s.get('sid'), parse_rfc3339(s.get('expires')))

  def extract_error(self, r):
    try:
      js = json.loads(r.data)
      if 'error' in js:
        self.error = js['error']
      elif 'message' in js:
        self.error = js['message']
    except BaseException:
      self.error = r.status

  def last_error(self):
    return self.error

  def last_response(self):
    return self.res
  def add_conference_admin (self, admin_id, conference_id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if admin_id is None:
            raise MissingRequiredArgument('property admin_id must be provided')
        payload['admin_id'] = admin_id
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if admin_id is not None:
            payload['admin_id'] = admin_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        uri = '%s/v2/conference/admin/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def add_conference_credential (self, conference_id, data, type, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if data is None:
            raise MissingRequiredArgument('property data must be provided')
        payload['data'] = data
        if type is None:
            raise MissingRequiredArgument('property type must be provided')
        payload['type'] = type
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if data is not None:
            payload['data'] = data
        if type is not None:
            payload['type'] = type
        uri = '%s/v2/conference/credentials/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def add_conference_date (self, conference_id, date, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if date is None:
            raise MissingRequiredArgument('property date must be provided')
        payload['date'] = date
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if date is not None:
            payload['date'] = date
        uri = '%s/v2/conference/date/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def add_conference_series_admin (self, admin_id, series_id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if admin_id is None:
            raise MissingRequiredArgument('property admin_id must be provided')
        payload['admin_id'] = admin_id
        if series_id is None:
            raise MissingRequiredArgument('property series_id must be provided')
        payload['series_id'] = series_id
        if admin_id is not None:
            payload['admin_id'] = admin_id
        if series_id is not None:
            payload['series_id'] = series_id
        uri = '%s/v2/conference_series/admin/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def add_conference_staff (self, conference_id, staff_id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if staff_id is None:
            raise MissingRequiredArgument('property staff_id must be provided')
        payload['staff_id'] = staff_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if staff_id is not None:
            payload['staff_id'] = staff_id
        uri = '%s/v2/conference/staff/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def add_conference_venue (self, conference_id, venue_id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if venue_id is None:
            raise MissingRequiredArgument('property venue_id must be provided')
        payload['venue_id'] = venue_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if venue_id is not None:
            payload['venue_id'] = venue_id
        uri = '%s/v2/conference/venue/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def add_featured_speaker (self, conference_id, description, display_name, avatar_url=None, speaker_id=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if description is None:
            raise MissingRequiredArgument('property description must be provided')
        payload['description'] = description
        if display_name is None:
            raise MissingRequiredArgument('property display_name must be provided')
        payload['display_name'] = display_name
        if avatar_url is not None:
            payload['avatar_url'] = avatar_url
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if description is not None:
            payload['description'] = description
        if display_name is not None:
            payload['display_name'] = display_name
        if speaker_id is not None:
            payload['speaker_id'] = speaker_id
        patterns = [re.compile('description#[a-z]+'), re.compile('display_name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/featured_speaker/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def add_session_type (self, abstract, conference_id, duration, name, submission_end=None, submission_start=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if abstract is None:
            raise MissingRequiredArgument('property abstract must be provided')
        payload['abstract'] = abstract
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if duration is None:
            raise MissingRequiredArgument('property duration must be provided')
        payload['duration'] = duration
        if name is None:
            raise MissingRequiredArgument('property name must be provided')
        payload['name'] = name
        if abstract is not None:
            payload['abstract'] = abstract
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if duration is not None:
            payload['duration'] = duration
        if name is not None:
            payload['name'] = name
        if submission_end is not None:
            payload['submission_end'] = submission_end
        if submission_start is not None:
            payload['submission_start'] = submission_start
        patterns = [re.compile('abstract#[a-z]+'), re.compile('name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/conference/session_type/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def add_sponsor (self, conference_id, group_name, name, url, logo_url=None, sort_order=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if group_name is None:
            raise MissingRequiredArgument('property group_name must be provided')
        payload['group_name'] = group_name
        if name is None:
            raise MissingRequiredArgument('property name must be provided')
        payload['name'] = name
        if url is None:
            raise MissingRequiredArgument('property url must be provided')
        payload['url'] = url
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if group_name is not None:
            payload['group_name'] = group_name
        if logo_url is not None:
            payload['logo_url'] = logo_url
        if name is not None:
            payload['name'] = name
        if sort_order is not None:
            payload['sort_order'] = sort_order
        if url is not None:
            payload['url'] = url
        patterns = [re.compile('name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/sponsor/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def confirm_temporary_email (self, confirmation_key, target_id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if confirmation_key is None:
            raise MissingRequiredArgument('property confirmation_key must be provided')
        payload['confirmation_key'] = confirmation_key
        if target_id is None:
            raise MissingRequiredArgument('property target_id must be provided')
        payload['target_id'] = target_id
        if confirmation_key is not None:
            payload['confirmation_key'] = confirmation_key
        if target_id is not None:
            payload['target_id'] = target_id
        uri = '%s/v2/email/confirm' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def create_blog_entry (self, conference_id, url, status=None, title=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if url is None:
            raise MissingRequiredArgument('property url must be provided')
        payload['url'] = url
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if status is not None:
            payload['status'] = status
        if title is not None:
            payload['title'] = title
        if url is not None:
            payload['url'] = url
        uri = '%s/v2/blog_entry/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def create_client_session (self, access_token, auth_via, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if access_token is None:
            raise MissingRequiredArgument('property access_token must be provided')
        payload['access_token'] = access_token
        if auth_via is None:
            raise MissingRequiredArgument('property auth_via must be provided')
        payload['auth_via'] = auth_via
        if access_token is not None:
            payload['access_token'] = access_token
        if auth_via is not None:
            payload['auth_via'] = auth_via
        uri = '%s/v2/client/session' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def create_conference (self, series_id, slug, title, cfp_lead_text=None, cfp_post_submit_instructions=None, cfp_pre_submit_instructions=None, contact_information=None, description=None, sub_title=None, timezone=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if series_id is None:
            raise MissingRequiredArgument('property series_id must be provided')
        payload['series_id'] = series_id
        if slug is None:
            raise MissingRequiredArgument('property slug must be provided')
        payload['slug'] = slug
        if title is None:
            raise MissingRequiredArgument('property title must be provided')
        payload['title'] = title
        if cfp_lead_text is not None:
            payload['cfp_lead_text'] = cfp_lead_text
        if cfp_post_submit_instructions is not None:
            payload['cfp_post_submit_instructions'] = cfp_post_submit_instructions
        if cfp_pre_submit_instructions is not None:
            payload['cfp_pre_submit_instructions'] = cfp_pre_submit_instructions
        if contact_information is not None:
            payload['contact_information'] = contact_information
        if description is not None:
            payload['description'] = description
        if series_id is not None:
            payload['series_id'] = series_id
        if slug is not None:
            payload['slug'] = slug
        if sub_title is not None:
            payload['sub_title'] = sub_title
        if timezone is not None:
            payload['timezone'] = timezone
        if title is not None:
            payload['title'] = title
        patterns = [re.compile('cfp_lead_text#[a-z]+'), re.compile('cfp_post_submit_instructions#[a-z]+'), re.compile('cfp_pre_submit_instructions#[a-z]+'), re.compile('contact_information#[a-z]+'), re.compile('description#[a-z]+'), re.compile('sub_title#[a-z]+'), re.compile('title#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/conference/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def create_conference_series (self, session, slug, title, description=None, sid=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if session is None:
            raise MissingRequiredArgument('property session must be provided')
        payload['session'] = session
        if slug is None:
            raise MissingRequiredArgument('property slug must be provided')
        payload['slug'] = slug
        if title is None:
            raise MissingRequiredArgument('property title must be provided')
        payload['title'] = title
        if description is not None:
            payload['description'] = description
        if sid is not None:
            payload['sid'] = sid
        if slug is not None:
            payload['slug'] = slug
        if title is not None:
            payload['title'] = title
        uri = '%s/v2/conference_series/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def create_external_resource (self, conference_id, title, url, description=None, image_url=None, sort_order=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if title is None:
            raise MissingRequiredArgument('property title must be provided')
        payload['title'] = title
        if url is None:
            raise MissingRequiredArgument('property url must be provided')
        payload['url'] = url
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if description is not None:
            payload['description'] = description
        if image_url is not None:
            payload['image_url'] = image_url
        if sort_order is not None:
            payload['sort_order'] = sort_order
        if title is not None:
            payload['title'] = title
        if url is not None:
            payload['url'] = url
        patterns = [re.compile('description#[a-z]+'), re.compile('title#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/external_resource/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def create_question (self, body, session_id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if body is None:
            raise MissingRequiredArgument('property body must be provided')
        payload['body'] = body
        if session_id is None:
            raise MissingRequiredArgument('property session_id must be provided')
        payload['session_id'] = session_id
        if body is not None:
            payload['body'] = body
        if session_id is not None:
            payload['session_id'] = session_id
        uri = '%s/v2/question/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def create_room (self, name, venue_id, capacity=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if name is None:
            raise MissingRequiredArgument('property name must be provided')
        payload['name'] = name
        if venue_id is None:
            raise MissingRequiredArgument('property venue_id must be provided')
        payload['venue_id'] = venue_id
        if capacity is not None:
            payload['capacity'] = capacity
        if name is not None:
            payload['name'] = name
        if venue_id is not None:
            payload['venue_id'] = venue_id
        patterns = [re.compile('name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/room/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def create_session (self, conference_id, session_type_id, speaker_id, abstract=None, category=None, material_level=None, materials_release=None, memo=None, photo_release=None, recording_release=None, slide_language=None, slide_subtitles=None, slide_url=None, spoken_language=None, tags=None, title=None, video_url=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if session_type_id is None:
            raise MissingRequiredArgument('property session_type_id must be provided')
        payload['session_type_id'] = session_type_id
        if speaker_id is None:
            raise MissingRequiredArgument('property speaker_id must be provided')
        payload['speaker_id'] = speaker_id
        if abstract is not None:
            payload['abstract'] = abstract
        if category is not None:
            payload['category'] = category
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if material_level is not None:
            payload['material_level'] = material_level
        if materials_release is not None:
            payload['materials_release'] = materials_release
        if memo is not None:
            payload['memo'] = memo
        if photo_release is not None:
            payload['photo_release'] = photo_release
        if recording_release is not None:
            payload['recording_release'] = recording_release
        if session_type_id is not None:
            payload['session_type_id'] = session_type_id
        if slide_language is not None:
            payload['slide_language'] = slide_language
        if slide_subtitles is not None:
            payload['slide_subtitles'] = slide_subtitles
        if slide_url is not None:
            payload['slide_url'] = slide_url
        if speaker_id is not None:
            payload['speaker_id'] = speaker_id
        if spoken_language is not None:
            payload['spoken_language'] = spoken_language
        if tags is not None:
            payload['tags'] = tags
        if title is not None:
            payload['title'] = title
        if video_url is not None:
            payload['video_url'] = video_url
        patterns = [re.compile('abstract#[a-z]+'), re.compile('title#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/session/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def create_session_survey_response (self, material_quality, overall_rating, session_id, speaker_knowledge, speaker_presentation, user_prior_knowledge, comment_good=None, comment_improvement=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if material_quality is None:
            raise MissingRequiredArgument('property material_quality must be provided')
        payload['material_quality'] = material_quality
        if overall_rating is None:
            raise MissingRequiredArgument('property overall_rating must be provided')
        payload['overall_rating'] = overall_rating
        if session_id is None:
            raise MissingRequiredArgument('property session_id must be provided')
        payload['session_id'] = session_id
        if speaker_knowledge is None:
            raise MissingRequiredArgument('property speaker_knowledge must be provided')
        payload['speaker_knowledge'] = speaker_knowledge
        if speaker_presentation is None:
            raise MissingRequiredArgument('property speaker_presentation must be provided')
        payload['speaker_presentation'] = speaker_presentation
        if user_prior_knowledge is None:
            raise MissingRequiredArgument('property user_prior_knowledge must be provided')
        payload['user_prior_knowledge'] = user_prior_knowledge
        if comment_good is not None:
            payload['comment_good'] = comment_good
        if comment_improvement is not None:
            payload['comment_improvement'] = comment_improvement
        if material_quality is not None:
            payload['material_quality'] = material_quality
        if overall_rating is not None:
            payload['overall_rating'] = overall_rating
        if session_id is not None:
            payload['session_id'] = session_id
        if speaker_knowledge is not None:
            payload['speaker_knowledge'] = speaker_knowledge
        if speaker_presentation is not None:
            payload['speaker_presentation'] = speaker_presentation
        if user_prior_knowledge is not None:
            payload['user_prior_knowledge'] = user_prior_knowledge
        uri = '%s/v2/survey_session_response/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def create_temporary_email (self, email, target_id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if email is None:
            raise MissingRequiredArgument('property email must be provided')
        payload['email'] = email
        if target_id is None:
            raise MissingRequiredArgument('property target_id must be provided')
        payload['target_id'] = target_id
        if email is not None:
            payload['email'] = email
        if target_id is not None:
            payload['target_id'] = target_id
        uri = '%s/v2/email/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def create_track (self, conference_id, room_id, name=None, sort_order=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if room_id is None:
            raise MissingRequiredArgument('property room_id must be provided')
        payload['room_id'] = room_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if name is not None:
            payload['name'] = name
        if room_id is not None:
            payload['room_id'] = room_id
        if sort_order is not None:
            payload['sort_order'] = sort_order
        patterns = [re.compile('name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/track/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def create_user (self, auth_user_id, auth_via, nickname, avatar_url=None, email=None, first_name=None, lang=None, last_name=None, tshirt_size=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if auth_user_id is None:
            raise MissingRequiredArgument('property auth_user_id must be provided')
        payload['auth_user_id'] = auth_user_id
        if auth_via is None:
            raise MissingRequiredArgument('property auth_via must be provided')
        payload['auth_via'] = auth_via
        if nickname is None:
            raise MissingRequiredArgument('property nickname must be provided')
        payload['nickname'] = nickname
        if auth_user_id is not None:
            payload['auth_user_id'] = auth_user_id
        if auth_via is not None:
            payload['auth_via'] = auth_via
        if avatar_url is not None:
            payload['avatar_url'] = avatar_url
        if email is not None:
            payload['email'] = email
        if first_name is not None:
            payload['first_name'] = first_name
        if lang is not None:
            payload['lang'] = lang
        if last_name is not None:
            payload['last_name'] = last_name
        if nickname is not None:
            payload['nickname'] = nickname
        if tshirt_size is not None:
            payload['tshirt_size'] = tshirt_size
        patterns = [re.compile('first_name#[a-z]+'), re.compile('last_name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/user/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def create_venue (self, address, name, latitude=None, longitude=None, place_id=None, url=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if address is None:
            raise MissingRequiredArgument('property address must be provided')
        payload['address'] = address
        if name is None:
            raise MissingRequiredArgument('property name must be provided')
        payload['name'] = name
        if address is not None:
            payload['address'] = address
        if latitude is not None:
            payload['latitude'] = latitude
        if longitude is not None:
            payload['longitude'] = longitude
        if name is not None:
            payload['name'] = name
        if place_id is not None:
            payload['place_id'] = place_id
        if url is not None:
            payload['url'] = url
        patterns = [re.compile('address#[a-z]+'), re.compile('name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/venue/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_blog_entry (self, id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v2/blog_entry/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_conference (self, id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v2/conference/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_conference_admin (self, admin_id, conference_id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if admin_id is None:
            raise MissingRequiredArgument('property admin_id must be provided')
        payload['admin_id'] = admin_id
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if admin_id is not None:
            payload['admin_id'] = admin_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        uri = '%s/v2/conference/admin/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_conference_date (self, conference_id, date, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if date is None:
            raise MissingRequiredArgument('property date must be provided')
        payload['date'] = date
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if date is not None:
            payload['date'] = date
        uri = '%s/v2/conference/date/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_conference_series (self, id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v2/conference_series/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_conference_staff (self, conference_id, staff_id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if staff_id is None:
            raise MissingRequiredArgument('property staff_id must be provided')
        payload['staff_id'] = staff_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if staff_id is not None:
            payload['staff_id'] = staff_id
        uri = '%s/v2/conference/staff/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_conference_venue (self, conference_id, venue_id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if venue_id is None:
            raise MissingRequiredArgument('property venue_id must be provided')
        payload['venue_id'] = venue_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if venue_id is not None:
            payload['venue_id'] = venue_id
        uri = '%s/v2/conference/venue/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_external_resource (self, id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v2/external_resource/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_featured_speaker (self, id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v2/featured_speaker/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_question (self, id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v2/question/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_room (self, id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v2/room/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_session (self, id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v2/session/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_session_type (self, id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v2/session_type/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_sponsor (self, id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v2/sponsor/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_track (self, id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v2/track/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_user (self, id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v2/user/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def delete_venue (self, id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v2/venue/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def get_conference_schedule (self, conference_id, lang=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if lang is not None:
            payload['lang'] = lang
        uri = '%s/v2/conference/schedule.ics' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def health_check (self, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        uri = '%s/' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def list_blog_entries (self, conference_id=None, lang=None, status=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if lang is not None:
            payload['lang'] = lang
        if status is not None:
            payload['status'] = status
        uri = '%s/v2/blog_entry/list' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def list_conference (self, lang=None, limit=None, organizers=None, range_end=None, range_start=None, since=None, status=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if lang is not None:
            payload['lang'] = lang
        if limit is not None:
            payload['limit'] = limit
        if organizers is not None:
            payload['organizers'] = organizers
        if range_end is not None:
            payload['range_end'] = range_end
        if range_start is not None:
            payload['range_start'] = range_start
        if since is not None:
            payload['since'] = since
        if status is not None:
            payload['status'] = status
        uri = '%s/v2/conference/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def list_conference_admin (self, conference_id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        uri = '%s/v2/conference/admin/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def list_conference_date (self, conference_id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        uri = '%s/v2/conference/date/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def list_conference_series (self, limit=None, since=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if limit is not None:
            payload['limit'] = limit
        if since is not None:
            payload['since'] = since
        uri = '%s/v2/conference_series/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def list_conference_staff (self, conference_id=None, lang=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if lang is not None:
            payload['lang'] = lang
        uri = '%s/v2/conference/staff/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def list_conferences_by_organizer (self, lang=None, limit=None, organizer_id=None, since=None, status=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if lang is not None:
            payload['lang'] = lang
        if limit is not None:
            payload['limit'] = limit
        if organizer_id is not None:
            payload['organizer_id'] = organizer_id
        if since is not None:
            payload['since'] = since
        if status is not None:
            payload['status'] = status
        uri = '%s/v2/conference/list_by_organizer' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def list_external_resource (self, conference_id, lang=None, limit=None, since=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if lang is not None:
            payload['lang'] = lang
        if limit is not None:
            payload['limit'] = limit
        if since is not None:
            payload['since'] = since
        uri = '%s/v2/external_resource/list' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def list_featured_speakers (self, conference_id=None, lang=None, limit=None, since=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if lang is not None:
            payload['lang'] = lang
        if limit is not None:
            payload['limit'] = limit
        if since is not None:
            payload['since'] = since
        uri = '%s/v2/featured_speaker/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def list_question (self, session_id, limit=None, since=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if session_id is None:
            raise MissingRequiredArgument('property session_id must be provided')
        payload['session_id'] = session_id
        if limit is not None:
            payload['limit'] = limit
        if session_id is not None:
            payload['session_id'] = session_id
        if since is not None:
            payload['since'] = since
        uri = '%s/v2/question/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def list_room (self, venue_id, lang=None, limit=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if venue_id is None:
            raise MissingRequiredArgument('property venue_id must be provided')
        payload['venue_id'] = venue_id
        if lang is not None:
            payload['lang'] = lang
        if limit is not None:
            payload['limit'] = limit
        if venue_id is not None:
            payload['venue_id'] = venue_id
        uri = '%s/v2/room/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def list_session_types_by_conference (self, conference_id=None, lang=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if lang is not None:
            payload['lang'] = lang
        uri = '%s/v2/session_type/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def list_sessions (self, conference_id=None, confirmed=None, lang=None, limit=None, range_end=None, range_start=None, since=None, speaker_id=None, status=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if confirmed is not None:
            payload['confirmed'] = confirmed
        if lang is not None:
            payload['lang'] = lang
        if limit is not None:
            payload['limit'] = limit
        if range_end is not None:
            payload['range_end'] = range_end
        if range_start is not None:
            payload['range_start'] = range_start
        if since is not None:
            payload['since'] = since
        if speaker_id is not None:
            payload['speaker_id'] = speaker_id
        if status is not None:
            payload['status'] = status
        uri = '%s/v2/session/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def list_sponsors (self, conference_id=None, lang=None, limit=None, since=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if lang is not None:
            payload['lang'] = lang
        if limit is not None:
            payload['limit'] = limit
        if since is not None:
            payload['since'] = since
        uri = '%s/v2/sponsor/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def list_user (self, lang=None, limit=None, pattern=None, since=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if lang is not None:
            payload['lang'] = lang
        if limit is not None:
            payload['limit'] = limit
        if pattern is not None:
            payload['pattern'] = pattern
        if since is not None:
            payload['since'] = since
        uri = '%s/v2/user/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def list_venue (self, lang=None, limit=None, since=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if lang is not None:
            payload['lang'] = lang
        if limit is not None:
            payload['limit'] = limit
        if since is not None:
            payload['since'] = since
        uri = '%s/v2/venue/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def lookup_blog_entry (self, id, lang=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        if lang is not None:
            payload['lang'] = lang
        uri = '%s/v2/blog_entry/lookup' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def lookup_conference (self, id, lang=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        if lang is not None:
            payload['lang'] = lang
        uri = '%s/v2/conference/lookup' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def lookup_conference_by_slug (self, slug, lang=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if slug is None:
            raise MissingRequiredArgument('property slug must be provided')
        payload['slug'] = slug
        if lang is not None:
            payload['lang'] = lang
        if slug is not None:
            payload['slug'] = slug
        uri = '%s/v2/conference/lookup_by_slug' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def lookup_conference_series (self, id, lang=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        if lang is not None:
            payload['lang'] = lang
        uri = '%s/v2/conference_series/lookup' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def lookup_external_resource (self, id, lang=None, limit=None, since=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        if lang is not None:
            payload['lang'] = lang
        if limit is not None:
            payload['limit'] = limit
        if since is not None:
            payload['since'] = since
        uri = '%s/v2/external_resource/lookup' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def lookup_featured_speaker (self, id, lang=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        if lang is not None:
            payload['lang'] = lang
        uri = '%s/v2/featured_speaker/lookup' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def lookup_room (self, id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v2/room/lookup' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def lookup_session (self, id, lang=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        if lang is not None:
            payload['lang'] = lang
        uri = '%s/v2/session/lookup' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def lookup_session_type (self, id, lang=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        if lang is not None:
            payload['lang'] = lang
        uri = '%s/v2/session_type/lookup' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def lookup_sponsor (self, id, lang=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        if lang is not None:
            payload['lang'] = lang
        uri = '%s/v2/sponsor/lookup' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def lookup_track (self, id, lang=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        if lang is not None:
            payload['lang'] = lang
        uri = '%s/v2/track/lookup' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def lookup_user (self, id, sid=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        if sid is not None:
            payload['sid'] = sid
        uri = '%s/v2/user/lookup' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def lookup_user_by_auth_user_id (self, auth_user_id, auth_via, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if auth_user_id is None:
            raise MissingRequiredArgument('property auth_user_id must be provided')
        payload['auth_user_id'] = auth_user_id
        if auth_via is None:
            raise MissingRequiredArgument('property auth_via must be provided')
        payload['auth_via'] = auth_via
        if auth_user_id is not None:
            payload['auth_user_id'] = auth_user_id
        if auth_via is not None:
            payload['auth_via'] = auth_via
        uri = '%s/v2/user/lookup_user_by_auth_user_id' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def lookup_venue (self, id, lang=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        if lang is not None:
            payload['lang'] = lang
        uri = '%s/v2/venue/lookup' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('GET', '%s?%s' % (uri, qs), headers=hdrs)
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def send_all_selection_result_notification (self, conference_id, force=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if force is not None:
            payload['force'] = force
        uri = '%s/v2/session/send_all_selection_result_notification' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def send_selection_result_notification (self, id, force=None, session_id=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if force is not None:
            payload['force'] = force
        if session_id is not None:
            payload['session_id'] = session_id
        uri = '%s/v2/session/send_selection_result_notification' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def set_session_video_cover (self, id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v2/session/video_cover' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def tweet_as_conference (self, conference_id, tweet, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if tweet is None:
            raise MissingRequiredArgument('property tweet must be provided')
        payload['tweet'] = tweet
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if tweet is not None:
            payload['tweet'] = tweet
        uri = '%s/v2/conference/tweet' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def update_blog_entry (self, id, status=None, title=None, url=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        if status is not None:
            payload['status'] = status
        if title is not None:
            payload['title'] = title
        if url is not None:
            payload['url'] = url
        uri = '%s/v2/blog_entry/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def update_conference (self, id, cfp_lead_text=None, cfp_post_submit_instructions=None, cfp_pre_submit_instructions=None, contact_information=None, description=None, redirect_url=None, slug=None, status=None, sub_title=None, timetable_available=None, timezone=None, title=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if cfp_lead_text is not None:
            payload['cfp_lead_text'] = cfp_lead_text
        if cfp_post_submit_instructions is not None:
            payload['cfp_post_submit_instructions'] = cfp_post_submit_instructions
        if cfp_pre_submit_instructions is not None:
            payload['cfp_pre_submit_instructions'] = cfp_pre_submit_instructions
        if contact_information is not None:
            payload['contact_information'] = contact_information
        if description is not None:
            payload['description'] = description
        if id is not None:
            payload['id'] = id
        if redirect_url is not None:
            payload['redirect_url'] = redirect_url
        if slug is not None:
            payload['slug'] = slug
        if status is not None:
            payload['status'] = status
        if sub_title is not None:
            payload['sub_title'] = sub_title
        if timetable_available is not None:
            payload['timetable_available'] = timetable_available
        if timezone is not None:
            payload['timezone'] = timezone
        if title is not None:
            payload['title'] = title
        patterns = [re.compile('cfp_lead_text#[a-z]+'), re.compile('cfp_post_submit_instructions#[a-z]+'), re.compile('cfp_pre_submit_instructions#[a-z]+'), re.compile('contact_information#[a-z]+'), re.compile('description#[a-z]+'), re.compile('sub_title#[a-z]+'), re.compile('title#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/conference/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def update_external_resource (self, id, description=None, image_url=None, sort_order=None, title=None, url=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if description is not None:
            payload['description'] = description
        if id is not None:
            payload['id'] = id
        if image_url is not None:
            payload['image_url'] = image_url
        if sort_order is not None:
            payload['sort_order'] = sort_order
        if title is not None:
            payload['title'] = title
        if url is not None:
            payload['url'] = url
        patterns = [re.compile('description#[a-z]+'), re.compile('title#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/external_resource/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return json.loads(res.data)
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def update_featured_speaker (self, id, avatar_url=None, description=None, display_name=None, speaker_id=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if avatar_url is not None:
            payload['avatar_url'] = avatar_url
        if description is not None:
            payload['description'] = description
        if display_name is not None:
            payload['display_name'] = display_name
        if id is not None:
            payload['id'] = id
        if speaker_id is not None:
            payload['speaker_id'] = speaker_id
        patterns = [re.compile('description#[a-z]+'), re.compile('display_name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/featured_speaker/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def update_room (self, id, capacity=None, name=None, venue_id=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if capacity is not None:
            payload['capacity'] = capacity
        if id is not None:
            payload['id'] = id
        if name is not None:
            payload['name'] = name
        if venue_id is not None:
            payload['venue_id'] = venue_id
        patterns = [re.compile('name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/room/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def update_session (self, id, abstract=None, category=None, conference_id=None, confirmed=None, duration=None, has_interpretation=None, material_level=None, materials_release=None, memo=None, photo_release=None, recording_release=None, session_type_id=None, slide_language=None, slide_subtitles=None, slide_url=None, sort_order=None, speaker_id=None, spoken_language=None, starts_on=None, status=None, tags=None, title=None, video_url=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if abstract is not None:
            payload['abstract'] = abstract
        if category is not None:
            payload['category'] = category
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if confirmed is not None:
            payload['confirmed'] = confirmed
        if duration is not None:
            payload['duration'] = duration
        if has_interpretation is not None:
            payload['has_interpretation'] = has_interpretation
        if id is not None:
            payload['id'] = id
        if material_level is not None:
            payload['material_level'] = material_level
        if materials_release is not None:
            payload['materials_release'] = materials_release
        if memo is not None:
            payload['memo'] = memo
        if photo_release is not None:
            payload['photo_release'] = photo_release
        if recording_release is not None:
            payload['recording_release'] = recording_release
        if session_type_id is not None:
            payload['session_type_id'] = session_type_id
        if slide_language is not None:
            payload['slide_language'] = slide_language
        if slide_subtitles is not None:
            payload['slide_subtitles'] = slide_subtitles
        if slide_url is not None:
            payload['slide_url'] = slide_url
        if sort_order is not None:
            payload['sort_order'] = sort_order
        if speaker_id is not None:
            payload['speaker_id'] = speaker_id
        if spoken_language is not None:
            payload['spoken_language'] = spoken_language
        if starts_on is not None:
            payload['starts_on'] = starts_on
        if status is not None:
            payload['status'] = status
        if tags is not None:
            payload['tags'] = tags
        if title is not None:
            payload['title'] = title
        if video_url is not None:
            payload['video_url'] = video_url
        patterns = [re.compile('abstract#[a-z]+'), re.compile('title#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/session/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def update_session_type (self, id, abstract=None, duration=None, is_default=None, name=None, submission_end=None, submission_start=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if abstract is not None:
            payload['abstract'] = abstract
        if duration is not None:
            payload['duration'] = duration
        if id is not None:
            payload['id'] = id
        if is_default is not None:
            payload['is_default'] = is_default
        if name is not None:
            payload['name'] = name
        if submission_end is not None:
            payload['submission_end'] = submission_end
        if submission_start is not None:
            payload['submission_start'] = submission_start
        patterns = [re.compile('abstract#[a-z]+'), re.compile('name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/session_type/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def update_sponsor (self, id, group_name=None, logo_url=None, name=None, sort_order=None, url=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if group_name is not None:
            payload['group_name'] = group_name
        if id is not None:
            payload['id'] = id
        if logo_url is not None:
            payload['logo_url'] = logo_url
        if name is not None:
            payload['name'] = name
        if sort_order is not None:
            payload['sort_order'] = sort_order
        if url is not None:
            payload['url'] = url
        patterns = [re.compile('name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/sponsor/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def update_track (self, id, name=None, room_id=None, sort_order=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        if name is not None:
            payload['name'] = name
        if room_id is not None:
            payload['room_id'] = room_id
        if sort_order is not None:
            payload['sort_order'] = sort_order
        patterns = [re.compile('name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/track/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def update_user (self, id, email=None, first_name=None, lang=None, last_name=None, nickname=None, tshirt_size=None, extra_headers=None, **args):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if email is not None:
            payload['email'] = email
        if first_name is not None:
            payload['first_name'] = first_name
        if id is not None:
            payload['id'] = id
        if lang is not None:
            payload['lang'] = lang
        if last_name is not None:
            payload['last_name'] = last_name
        if nickname is not None:
            payload['nickname'] = nickname
        if tshirt_size is not None:
            payload['tshirt_size'] = tshirt_size
        patterns = [re.compile('first_name#[a-z]+'), re.compile('last_name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v2/user/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def update_venue (self, id, latitude=None, longitude=None, name=None, place_id=None, url=None, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        if latitude is not None:
            payload['latitude'] = latitude
        if longitude is not None:
            payload['longitude'] = longitude
        if name is not None:
            payload['name'] = name
        if place_id is not None:
            payload['place_id'] = place_id
        if url is not None:
            payload['url'] = url
        uri = '%s/v2/venue/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


  def verify_user (self, id, extra_headers=None):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v2/user/verify' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
        if extra_headers:
            hdrs.update(extra_headers)
        res = self.http.request('POST', uri, headers=hdrs, body=json.dumps(payload))
        if self.debug:
            print(res)
        self.res = res
        if res.status != 200:
            self.extract_error(res)
            return None
        return True
    except BaseException as e:
        if self.debug:
            print("error during http access: " + repr(e))
        self.error = repr(e)
        return None


