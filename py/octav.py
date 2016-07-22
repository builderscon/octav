"""OCTAV Client Library"""
"""DO NOT EDIT: This file was generated from ../spec/v1/api.json on Fri Jul 22 13:49:34 2016"""

import json
import os
import requests

class Octav(object):
  def __init__(self, endpoint, key, secret, debug=False):
    if not endpoint:
      raise "endpoint is required"
    if not key:
      raise "key is required"
    if not secret:
      raise "secret is required"
    self.debug = debug
    self.endpoint = endpoint
    self.error = None
    self.key = key
    self.secret = secret
    self.session = requests.Session()
    self.session.mount('http://', requests.adapters.HTTPAdapter(max_retries=0))

  def extract_error(self, r):
    try:
      js = r.json()
      self.error = js["message"]
    except:
      self.error = r.status_code

  def last_error(self):
    return self.error

  def create_user (self, auth_user_id, auth_via, nickname, avatar_url=None, email=None, first_name=None, last_name=None, tshirt_size=None):
    payload = {}
    if auth_user_id is None:
            raise 'property auth_user_id must be provided'
    payload['auth_user_id'] = auth_user_id
    if auth_via is None:
            raise 'property auth_via must be provided'
    payload['auth_via'] = auth_via
    if nickname is None:
            raise 'property nickname must be provided'
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
    if last_name is not None:
        payload['last_name'] = last_name
    if nickname is not None:
        payload['nickname'] = nickname
    if tshirt_size is not None:
        payload['tshirt_size'] = tshirt_size
    uri = self.endpoint + "/user/create"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def lookup_user (self, id):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if id is not None:
        payload['id'] = id
    uri = self.endpoint + "/user/lookup"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def lookup_user_by_auth_user_id (self, auth_user_id, auth_via):
    payload = {}
    if auth_user_id is None:
            raise 'property auth_user_id must be provided'
    payload['auth_user_id'] = auth_user_id
    if auth_via is None:
            raise 'property auth_via must be provided'
    payload['auth_via'] = auth_via
    if auth_user_id is not None:
        payload['auth_user_id'] = auth_user_id
    if auth_via is not None:
        payload['auth_via'] = auth_via
    uri = self.endpoint + "/user/lookup_by_auth_user_id"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def update_user (self, id, user_id, email=None, first_name=None, last_name=None, nickname=None, tshirt_size=None):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if email is not None:
        payload['email'] = email
    if first_name is not None:
        payload['first_name'] = first_name
    if id is not None:
        payload['id'] = id
    if last_name is not None:
        payload['last_name'] = last_name
    if nickname is not None:
        payload['nickname'] = nickname
    if tshirt_size is not None:
        payload['tshirt_size'] = tshirt_size
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/user/update"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_user (self, id, user_id):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if id is not None:
        payload['id'] = id
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/user/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def list_user (self, lang=None, limit=None, since=None):
    payload = {}
    if lang is not None:
        payload['lang'] = lang
    if limit is not None:
        payload['limit'] = limit
    if since is not None:
        payload['since'] = since
    uri = self.endpoint + "/user/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def create_venue (self, address, name, user_id, latitude=None, longitude=None):
    payload = {}
    if address is None:
            raise 'property address must be provided'
    payload['address'] = address
    if name is None:
            raise 'property name must be provided'
    payload['name'] = name
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if address is not None:
        payload['address'] = address
    if latitude is not None:
        payload['latitude'] = latitude
    if longitude is not None:
        payload['longitude'] = longitude
    if name is not None:
        payload['name'] = name
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/venue/create"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def list_venue (self, lang=None, limit=None, since=None):
    payload = {}
    if lang is not None:
        payload['lang'] = lang
    if limit is not None:
        payload['limit'] = limit
    if since is not None:
        payload['since'] = since
    uri = self.endpoint + "/venue/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def lookup_venue (self, id):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if id is not None:
        payload['id'] = id
    uri = self.endpoint + "/venue/lookup"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def update_venue (self, id, user_id):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if id is not None:
        payload['id'] = id
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/venue/update"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_venue (self, id, user_id):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if id is not None:
        payload['id'] = id
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/venue/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def create_room (self, name, user_id, venue_id, capacity=None):
    payload = {}
    if name is None:
            raise 'property name must be provided'
    payload['name'] = name
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if venue_id is None:
            raise 'property venue_id must be provided'
    payload['venue_id'] = venue_id
    if capacity is not None:
        payload['capacity'] = capacity
    if name is not None:
        payload['name'] = name
    if user_id is not None:
        payload['user_id'] = user_id
    if venue_id is not None:
        payload['venue_id'] = venue_id
    uri = self.endpoint + "/room/create"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def update_room (self, id, user_id, capacity=None, name=None, venue_id=None):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if capacity is not None:
        payload['capacity'] = capacity
    if id is not None:
        payload['id'] = id
    if name is not None:
        payload['name'] = name
    if user_id is not None:
        payload['user_id'] = user_id
    if venue_id is not None:
        payload['venue_id'] = venue_id
    uri = self.endpoint + "/room/update"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def lookup_room (self, id):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if id is not None:
        payload['id'] = id
    uri = self.endpoint + "/room/lookup"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def delete_room (self, id, user_id):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if id is not None:
        payload['id'] = id
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/room/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def list_room (self, venue_id, lang=None, limit=None):
    payload = {}
    if venue_id is None:
            raise 'property venue_id must be provided'
    payload['venue_id'] = venue_id
    if lang is not None:
        payload['lang'] = lang
    if limit is not None:
        payload['limit'] = limit
    if venue_id is not None:
        payload['venue_id'] = venue_id
    uri = self.endpoint + "/room/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def create_conference_series (self, slug, title, user_id):
    payload = {}
    if slug is None:
            raise 'property slug must be provided'
    payload['slug'] = slug
    if title is None:
            raise 'property title must be provided'
    payload['title'] = title
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if slug is not None:
        payload['slug'] = slug
    if title is not None:
        payload['title'] = title
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/conference_series/create"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def list_conference_series (self, limit=None, since=None):
    payload = {}
    if limit is not None:
        payload['limit'] = limit
    if since is not None:
        payload['since'] = since
    uri = self.endpoint + "/conference_series/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def add_conference_series_admin (self, admin_id, series_id, user_id):
    payload = {}
    if admin_id is None:
            raise 'property admin_id must be provided'
    payload['admin_id'] = admin_id
    if series_id is None:
            raise 'property series_id must be provided'
    payload['series_id'] = series_id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if admin_id is not None:
        payload['admin_id'] = admin_id
    if series_id is not None:
        payload['series_id'] = series_id
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/conference_series/admin/add"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def create_conference (self, series_id, slug, title, user_id, description=None, sub_title=None):
    payload = {}
    if series_id is None:
            raise 'property series_id must be provided'
    payload['series_id'] = series_id
    if slug is None:
            raise 'property slug must be provided'
    payload['slug'] = slug
    if title is None:
            raise 'property title must be provided'
    payload['title'] = title
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if description is not None:
        payload['description'] = description
    if series_id is not None:
        payload['series_id'] = series_id
    if slug is not None:
        payload['slug'] = slug
    if sub_title is not None:
        payload['sub_title'] = sub_title
    if title is not None:
        payload['title'] = title
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/conference/create"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def add_conference_dates (self, conference_id, dates, user_id):
    payload = {}
    if conference_id is None:
            raise 'property conference_id must be provided'
    payload['conference_id'] = conference_id
    if dates is None:
            raise 'property dates must be provided'
    payload['dates'] = dates
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if conference_id is not None:
        payload['conference_id'] = conference_id
    if dates is not None:
        payload['dates'] = dates
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/conference/dates/add"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_conference_dates (self, conference_id, dates, user_id):
    payload = {}
    if conference_id is None:
            raise 'property conference_id must be provided'
    payload['conference_id'] = conference_id
    if dates is None:
            raise 'property dates must be provided'
    payload['dates'] = dates
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if conference_id is not None:
        payload['conference_id'] = conference_id
    if dates is not None:
        payload['dates'] = dates
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/conference/dates/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def add_conference_admin (self, admin_id, conference_id, user_id):
    payload = {}
    if admin_id is None:
            raise 'property admin_id must be provided'
    payload['admin_id'] = admin_id
    if conference_id is None:
            raise 'property conference_id must be provided'
    payload['conference_id'] = conference_id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if admin_id is not None:
        payload['admin_id'] = admin_id
    if conference_id is not None:
        payload['conference_id'] = conference_id
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/conference/admin/add"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_conference_admin (self, admin_id, conference_id, user_id):
    payload = {}
    if admin_id is None:
            raise 'property admin_id must be provided'
    payload['admin_id'] = admin_id
    if conference_id is None:
            raise 'property conference_id must be provided'
    payload['conference_id'] = conference_id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if admin_id is not None:
        payload['admin_id'] = admin_id
    if conference_id is not None:
        payload['conference_id'] = conference_id
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/conference/admin/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def add_conference_venue (self, conference_id, user_id, venue_id):
    payload = {}
    if conference_id is None:
            raise 'property conference_id must be provided'
    payload['conference_id'] = conference_id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if venue_id is None:
            raise 'property venue_id must be provided'
    payload['venue_id'] = venue_id
    if conference_id is not None:
        payload['conference_id'] = conference_id
    if user_id is not None:
        payload['user_id'] = user_id
    if venue_id is not None:
        payload['venue_id'] = venue_id
    uri = self.endpoint + "/conference/venue/add"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_conference_venue (self, conference_id, user_id, venue_id):
    payload = {}
    if conference_id is None:
            raise 'property conference_id must be provided'
    payload['conference_id'] = conference_id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if venue_id is None:
            raise 'property venue_id must be provided'
    payload['venue_id'] = venue_id
    if conference_id is not None:
        payload['conference_id'] = conference_id
    if user_id is not None:
        payload['user_id'] = user_id
    if venue_id is not None:
        payload['venue_id'] = venue_id
    uri = self.endpoint + "/conference/venue/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def lookup_conference (self, id, lang=None):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if id is not None:
        payload['id'] = id
    if lang is not None:
        payload['lang'] = lang
    uri = self.endpoint + "/conference/lookup"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def lookup_conference_by_slug (self, slug, lang=None):
    payload = {}
    if slug is None:
            raise 'property slug must be provided'
    payload['slug'] = slug
    if lang is not None:
        payload['lang'] = lang
    if slug is not None:
        payload['slug'] = slug
    uri = self.endpoint + "/conference/lookup_by_slug"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def list_conference (self, lang=None, limit=None, range_end=None, range_start=None, since=None, status=None):
    payload = {}
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
    if status is not None:
        payload['status'] = status
    uri = self.endpoint + "/conference/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def update_conference (self, id, user_id, description=None, slug=None, status=None, sub_title=None, title=None):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if description is not None:
        payload['description'] = description
    if id is not None:
        payload['id'] = id
    if slug is not None:
        payload['slug'] = slug
    if status is not None:
        payload['status'] = status
    if sub_title is not None:
        payload['sub_title'] = sub_title
    if title is not None:
        payload['title'] = title
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/conference/update"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_conference_series (self, id, user_id):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if id is not None:
        payload['id'] = id
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/conference_series/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_conference (self, id):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if id is not None:
        payload['id'] = id
    uri = self.endpoint + "/conference/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def create_session (self, abstract, conference_id, duration, speaker_id, title, user_id, category=None, material_level=None, memo=None, photo_permission=None, slide_language=None, slide_subtitles=None, slide_url=None, spoken_language=None, tags=None, video_permission=None, video_url=None):
    payload = {}
    if abstract is None:
            raise 'property abstract must be provided'
    payload['abstract'] = abstract
    if conference_id is None:
            raise 'property conference_id must be provided'
    payload['conference_id'] = conference_id
    if duration is None:
            raise 'property duration must be provided'
    payload['duration'] = duration
    if speaker_id is None:
            raise 'property speaker_id must be provided'
    payload['speaker_id'] = speaker_id
    if title is None:
            raise 'property title must be provided'
    payload['title'] = title
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if abstract is not None:
        payload['abstract'] = abstract
    if category is not None:
        payload['category'] = category
    if conference_id is not None:
        payload['conference_id'] = conference_id
    if duration is not None:
        payload['duration'] = duration
    if material_level is not None:
        payload['material_level'] = material_level
    if memo is not None:
        payload['memo'] = memo
    if photo_permission is not None:
        payload['photo_permission'] = photo_permission
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
    if user_id is not None:
        payload['user_id'] = user_id
    if video_permission is not None:
        payload['video_permission'] = video_permission
    if video_url is not None:
        payload['video_url'] = video_url
    uri = self.endpoint + "/session/create"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def lookup_session (self, id):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if id is not None:
        payload['id'] = id
    uri = self.endpoint + "/session/lookup"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def delete_session (self, id, user_id):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if id is not None:
        payload['id'] = id
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/session/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def update_session (self, id, user_id, abstract=None, category=None, conference_id=None, confirmed=None, duration=None, has_interpretation=None, material_level=None, memo=None, photo_permission=None, slide_language=None, slide_subtitles=None, slide_url=None, sort_order=None, speaker_id=None, spoken_language=None, status=None, tags=None, title=None, video_permission=None, video_url=None):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
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
    if memo is not None:
        payload['memo'] = memo
    if photo_permission is not None:
        payload['photo_permission'] = photo_permission
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
    if status is not None:
        payload['status'] = status
    if tags is not None:
        payload['tags'] = tags
    if title is not None:
        payload['title'] = title
    if user_id is not None:
        payload['user_id'] = user_id
    if video_permission is not None:
        payload['video_permission'] = video_permission
    if video_url is not None:
        payload['video_url'] = video_url
    uri = self.endpoint + "/session/update"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def list_session_by_conference (self, conference_id, date=None):
    payload = {}
    if conference_id is None:
            raise 'property conference_id must be provided'
    payload['conference_id'] = conference_id
    if conference_id is not None:
        payload['conference_id'] = conference_id
    if date is not None:
        payload['date'] = date
    uri = self.endpoint + "/schedule/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def create_question (self, body, session_id, user_id):
    payload = {}
    if body is None:
            raise 'property body must be provided'
    payload['body'] = body
    if session_id is None:
            raise 'property session_id must be provided'
    payload['session_id'] = session_id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if body is not None:
        payload['body'] = body
    if session_id is not None:
        payload['session_id'] = session_id
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/question/create"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def delete_question (self, id):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if id is not None:
        payload['id'] = id
    uri = self.endpoint + "/question/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def list_question (self, session_id, limit=None, since=None):
    payload = {}
    if session_id is None:
            raise 'property session_id must be provided'
    payload['session_id'] = session_id
    if limit is not None:
        payload['limit'] = limit
    if session_id is not None:
        payload['session_id'] = session_id
    if since is not None:
        payload['since'] = since
    uri = self.endpoint + "/question/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def create_session_survey_response (self, material_quality, overall_rating, session_id, speaker_knowledge, speaker_presentation, user_id, user_prior_knowledge, comment_good=None, comment_improvement=None):
    payload = {}
    if material_quality is None:
            raise 'property material_quality must be provided'
    payload['material_quality'] = material_quality
    if overall_rating is None:
            raise 'property overall_rating must be provided'
    payload['overall_rating'] = overall_rating
    if session_id is None:
            raise 'property session_id must be provided'
    payload['session_id'] = session_id
    if speaker_knowledge is None:
            raise 'property speaker_knowledge must be provided'
    payload['speaker_knowledge'] = speaker_knowledge
    if speaker_presentation is None:
            raise 'property speaker_presentation must be provided'
    payload['speaker_presentation'] = speaker_presentation
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if user_prior_knowledge is None:
            raise 'property user_prior_knowledge must be provided'
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
    if user_id is not None:
        payload['user_id'] = user_id
    if user_prior_knowledge is not None:
        payload['user_prior_knowledge'] = user_prior_knowledge
    uri = self.endpoint + "/survey_session_response/create"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def add_featured_speaker (self, conference_id, description, display_name, avatar_url=None, speaker_id=None, user_id=None):
    payload = {}
    if conference_id is None:
            raise 'property conference_id must be provided'
    payload['conference_id'] = conference_id
    if description is None:
            raise 'property description must be provided'
    payload['description'] = description
    if display_name is None:
            raise 'property display_name must be provided'
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
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/featured_speaker/add"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def lookup_featured_speaker (self, id, lang=None):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if id is not None:
        payload['id'] = id
    if lang is not None:
        payload['lang'] = lang
    uri = self.endpoint + "/featured_speaker/lookup"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def list_featured_speakers (self, conference_id=None, lang=None, limit=None, since=None):
    payload = {}
    if conference_id is not None:
        payload['conference_id'] = conference_id
    if lang is not None:
        payload['lang'] = lang
    if limit is not None:
        payload['limit'] = limit
    if since is not None:
        payload['since'] = since
    uri = self.endpoint + "/featured_speaker/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def update_featured_speaker (self, id, user_id, avatar_url=None, description=None, display_name=None, speaker_id=None):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
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
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/featured_speaker/update"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_featured_speaker (self, id, user_id):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if id is not None:
        payload['id'] = id
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/featured_speaker/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def add_sponsor (self, conference_id, group_name, logo_url1, name, url, user_id, logo_url2=None, logo_url3=None, sort_order=None):
    payload = {}
    if conference_id is None:
            raise 'property conference_id must be provided'
    payload['conference_id'] = conference_id
    if group_name is None:
            raise 'property group_name must be provided'
    payload['group_name'] = group_name
    if logo_url1 is None:
            raise 'property logo_url1 must be provided'
    payload['logo_url1'] = logo_url1
    if name is None:
            raise 'property name must be provided'
    payload['name'] = name
    if url is None:
            raise 'property url must be provided'
    payload['url'] = url
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if conference_id is not None:
        payload['conference_id'] = conference_id
    if group_name is not None:
        payload['group_name'] = group_name
    if logo_url1 is not None:
        payload['logo_url1'] = logo_url1
    if logo_url2 is not None:
        payload['logo_url2'] = logo_url2
    if logo_url3 is not None:
        payload['logo_url3'] = logo_url3
    if name is not None:
        payload['name'] = name
    if sort_order is not None:
        payload['sort_order'] = sort_order
    if url is not None:
        payload['url'] = url
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/sponsor/add"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def lookup_sponsor (self, id, lang=None):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if id is not None:
        payload['id'] = id
    if lang is not None:
        payload['lang'] = lang
    uri = self.endpoint + "/sponsor/lookup"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def list_sponsors (self, conference_id=None, lang=None, limit=None, since=None):
    payload = {}
    if conference_id is not None:
        payload['conference_id'] = conference_id
    if lang is not None:
        payload['lang'] = lang
    if limit is not None:
        payload['limit'] = limit
    if since is not None:
        payload['since'] = since
    uri = self.endpoint + "/sponsor/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def update_sponsor (self, id, user_id, group_name=None, logo_url1=None, logo_url2=None, logo_url3=None, name=None, sort_order=None, url=None):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if group_name is not None:
        payload['group_name'] = group_name
    if id is not None:
        payload['id'] = id
    if logo_url1 is not None:
        payload['logo_url1'] = logo_url1
    if logo_url2 is not None:
        payload['logo_url2'] = logo_url2
    if logo_url3 is not None:
        payload['logo_url3'] = logo_url3
    if name is not None:
        payload['name'] = name
    if sort_order is not None:
        payload['sort_order'] = sort_order
    if url is not None:
        payload['url'] = url
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/sponsor/update"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_sponsor (self, id, user_id):
    payload = {}
    if id is None:
            raise 'property id must be provided'
    payload['id'] = id
    if user_id is None:
            raise 'property user_id must be provided'
    payload['user_id'] = user_id
    if id is not None:
        payload['id'] = id
    if user_id is not None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/sponsor/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

