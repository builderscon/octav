"""OCTAV Client Library"""
"""DO NOT EDIT: This file was generated from ../spec/v1/api.json on Fri Jul 15 08:22:33 2016"""

import json
import os
import requests

class Octav(object):
  def __init__(self, endpoint=None, key=None, secret=None, debug=False):
    if not endpoint:
      raise "endpoint is required"
    if not key:
      raise "key is required"
    if not secret:
      raise "secret is required"
    self.debug = debug
    self.endpoint = endpoint
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

  def create_user (self, tshirt_size=None, auth_user_id=None, first_name=None, email=None, last_name=None, nickname=None, avatar_url=None, auth_via=None):
    if nickname is None:
            raise "property \"" + required + "\" must be provided"
    if auth_via is None:
            raise "property \"" + required + "\" must be provided"
    if auth_user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not tshirt_size is None:
        payload['tshirt_size'] = tshirt_size
    if not auth_user_id is None:
        payload['auth_user_id'] = auth_user_id
    if not first_name is None:
        payload['first_name'] = first_name
    if not email is None:
        payload['email'] = email
    if not last_name is None:
        payload['last_name'] = last_name
    if not nickname is None:
        payload['nickname'] = nickname
    if not avatar_url is None:
        payload['avatar_url'] = avatar_url
    if not auth_via is None:
        payload['auth_via'] = auth_via
    uri = self.endpoint + "/user/create"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def lookup_user (self, id=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not id is None:
        payload['id'] = id
    uri = self.endpoint + "/user/lookup"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def lookup_user_by_auth_user_id (self, auth_via=None, auth_user_id=None):
    if auth_via is None:
            raise "property \"" + required + "\" must be provided"
    if auth_user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not auth_via is None:
        payload['auth_via'] = auth_via
    if not auth_user_id is None:
        payload['auth_user_id'] = auth_user_id
    uri = self.endpoint + "/user/lookup_by_auth_user_id"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def update_user (self, id=None, tshirt_size=None, last_name=None, user_id=None, nickname=None, first_name=None, email=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not id is None:
        payload['id'] = id
    if not tshirt_size is None:
        payload['tshirt_size'] = tshirt_size
    if not last_name is None:
        payload['last_name'] = last_name
    if not user_id is None:
        payload['user_id'] = user_id
    if not nickname is None:
        payload['nickname'] = nickname
    if not first_name is None:
        payload['first_name'] = first_name
    if not email is None:
        payload['email'] = email
    uri = self.endpoint + "/user/update"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_user (self, user_id=None, id=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not user_id is None:
        payload['user_id'] = user_id
    if not id is None:
        payload['id'] = id
    uri = self.endpoint + "/user/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def list_user (self, limit=None, lang=None, since=None):
    payload = {}
    if not limit is None:
        payload['limit'] = limit
    if not lang is None:
        payload['lang'] = lang
    if not since is None:
        payload['since'] = since
    uri = self.endpoint + "/user/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def create_venue (self, name=None, latitude=None, user_id=None, address=None, longitude=None):
    if name is None:
            raise "property \"" + required + "\" must be provided"
    if address is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not name is None:
        payload['name'] = name
    if not latitude is None:
        payload['latitude'] = latitude
    if not user_id is None:
        payload['user_id'] = user_id
    if not address is None:
        payload['address'] = address
    if not longitude is None:
        payload['longitude'] = longitude
    uri = self.endpoint + "/venue/create"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def list_venue (self, since=None, lang=None, limit=None):
    payload = {}
    if not since is None:
        payload['since'] = since
    if not lang is None:
        payload['lang'] = lang
    if not limit is None:
        payload['limit'] = limit
    uri = self.endpoint + "/venue/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def lookup_venue (self, id=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not id is None:
        payload['id'] = id
    uri = self.endpoint + "/venue/lookup"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def update_venue (self, user_id=None, id=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not user_id is None:
        payload['user_id'] = user_id
    if not id is None:
        payload['id'] = id
    uri = self.endpoint + "/venue/update"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_venue (self, user_id=None, id=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not user_id is None:
        payload['user_id'] = user_id
    if not id is None:
        payload['id'] = id
    uri = self.endpoint + "/venue/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def create_room (self, capacity=None, user_id=None, venue_id=None, name=None):
    if venue_id is None:
            raise "property \"" + required + "\" must be provided"
    if name is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not capacity is None:
        payload['capacity'] = capacity
    if not user_id is None:
        payload['user_id'] = user_id
    if not venue_id is None:
        payload['venue_id'] = venue_id
    if not name is None:
        payload['name'] = name
    uri = self.endpoint + "/room/create"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def update_room (self, capacity=None, user_id=None, name=None, venue_id=None, id=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not capacity is None:
        payload['capacity'] = capacity
    if not user_id is None:
        payload['user_id'] = user_id
    if not name is None:
        payload['name'] = name
    if not venue_id is None:
        payload['venue_id'] = venue_id
    if not id is None:
        payload['id'] = id
    uri = self.endpoint + "/room/update"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def lookup_room (self, id=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not id is None:
        payload['id'] = id
    uri = self.endpoint + "/room/lookup"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def delete_room (self, id=None, user_id=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not id is None:
        payload['id'] = id
    if not user_id is None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/room/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def list_room (self, limit=None, venue_id=None, lang=None):
    if venue_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not limit is None:
        payload['limit'] = limit
    if not venue_id is None:
        payload['venue_id'] = venue_id
    if not lang is None:
        payload['lang'] = lang
    uri = self.endpoint + "/room/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def create_conference_series (self, slug=None, user_id=None):
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    if slug is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not slug is None:
        payload['slug'] = slug
    if not user_id is None:
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
    if not limit is None:
        payload['limit'] = limit
    if not since is None:
        payload['since'] = since
    uri = self.endpoint + "/conference_series/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def add_conference_series_admin (self, user_id=None, admin_id=None, series_id=None):
    if series_id is None:
            raise "property \"" + required + "\" must be provided"
    if admin_id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not user_id is None:
        payload['user_id'] = user_id
    if not admin_id is None:
        payload['admin_id'] = admin_id
    if not series_id is None:
        payload['series_id'] = series_id
    uri = self.endpoint + "/conference_series/admin/add"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def create_conference (self, user_id=None, sub_title=None, slug=None, series_id=None, description=None, title=None):
    if series_id is None:
            raise "property \"" + required + "\" must be provided"
    if title is None:
            raise "property \"" + required + "\" must be provided"
    if slug is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not user_id is None:
        payload['user_id'] = user_id
    if not sub_title is None:
        payload['sub_title'] = sub_title
    if not slug is None:
        payload['slug'] = slug
    if not series_id is None:
        payload['series_id'] = series_id
    if not description is None:
        payload['description'] = description
    if not title is None:
        payload['title'] = title
    uri = self.endpoint + "/conference/create"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def add_conference_dates (self, user_id=None, conference_id=None, dates=None):
    if conference_id is None:
            raise "property \"" + required + "\" must be provided"
    if dates is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not user_id is None:
        payload['user_id'] = user_id
    if not conference_id is None:
        payload['conference_id'] = conference_id
    if not dates is None:
        payload['dates'] = dates
    uri = self.endpoint + "/conference/dates/add"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_conference_dates (self, dates=None, conference_id=None, user_id=None):
    if conference_id is None:
            raise "property \"" + required + "\" must be provided"
    if dates is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not dates is None:
        payload['dates'] = dates
    if not conference_id is None:
        payload['conference_id'] = conference_id
    if not user_id is None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/conference/dates/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def add_conference_admin (self, conference_id=None, admin_id=None, user_id=None):
    if conference_id is None:
            raise "property \"" + required + "\" must be provided"
    if admin_id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not conference_id is None:
        payload['conference_id'] = conference_id
    if not admin_id is None:
        payload['admin_id'] = admin_id
    if not user_id is None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/conference/admin/add"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_conference_admin (self, conference_id=None, admin_id=None, user_id=None):
    if conference_id is None:
            raise "property \"" + required + "\" must be provided"
    if admin_id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not conference_id is None:
        payload['conference_id'] = conference_id
    if not admin_id is None:
        payload['admin_id'] = admin_id
    if not user_id is None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/conference/admin/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def add_conference_venue (self, conference_id=None, venue_id=None, user_id=None):
    if conference_id is None:
            raise "property \"" + required + "\" must be provided"
    if venue_id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not conference_id is None:
        payload['conference_id'] = conference_id
    if not venue_id is None:
        payload['venue_id'] = venue_id
    if not user_id is None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/conference/venue/add"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_conference_venue (self, venue_id=None, conference_id=None, user_id=None):
    if conference_id is None:
            raise "property \"" + required + "\" must be provided"
    if venue_id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not venue_id is None:
        payload['venue_id'] = venue_id
    if not conference_id is None:
        payload['conference_id'] = conference_id
    if not user_id is None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/conference/venue/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def lookup_conference (self, id=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not id is None:
        payload['id'] = id
    uri = self.endpoint + "/conference/lookup"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def lookup_conference_by_slug (self, slug=None, lang=None):
    if slug is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not slug is None:
        payload['slug'] = slug
    if not lang is None:
        payload['lang'] = lang
    uri = self.endpoint + "/conference/lookup_by_slug"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def list_conference (self, range_end=None, limit=None, range_start=None, since=None, lang=None):
    payload = {}
    if not range_end is None:
        payload['range_end'] = range_end
    if not limit is None:
        payload['limit'] = limit
    if not range_start is None:
        payload['range_start'] = range_start
    if not since is None:
        payload['since'] = since
    if not lang is None:
        payload['lang'] = lang
    uri = self.endpoint + "/conference/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def update_conference (self, slug=None, sub_title=None, title=None, description=None, id=None, user_id=None, starts_on=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not slug is None:
        payload['slug'] = slug
    if not sub_title is None:
        payload['sub_title'] = sub_title
    if not title is None:
        payload['title'] = title
    if not description is None:
        payload['description'] = description
    if not id is None:
        payload['id'] = id
    if not user_id is None:
        payload['user_id'] = user_id
    if not starts_on is None:
        payload['starts_on'] = starts_on
    uri = self.endpoint + "/conference/update"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_conference_series (self, user_id=None, id=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not user_id is None:
        payload['user_id'] = user_id
    if not id is None:
        payload['id'] = id
    uri = self.endpoint + "/conference_series/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_conference (self, id=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not id is None:
        payload['id'] = id
    uri = self.endpoint + "/conference/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def create_session (self, tags=None, slide_url=None, video_url=None, speaker_id=None, abstract=None, video_permission=None, user_id=None, slide_language=None, photo_permission=None, title=None, duration=None, spoken_language=None, memo=None, category=None, conference_id=None, material_level=None, slide_subtitles=None):
    if conference_id is None:
            raise "property \"" + required + "\" must be provided"
    if speaker_id is None:
            raise "property \"" + required + "\" must be provided"
    if title is None:
            raise "property \"" + required + "\" must be provided"
    if abstract is None:
            raise "property \"" + required + "\" must be provided"
    if duration is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not tags is None:
        payload['tags'] = tags
    if not slide_url is None:
        payload['slide_url'] = slide_url
    if not video_url is None:
        payload['video_url'] = video_url
    if not speaker_id is None:
        payload['speaker_id'] = speaker_id
    if not abstract is None:
        payload['abstract'] = abstract
    if not video_permission is None:
        payload['video_permission'] = video_permission
    if not user_id is None:
        payload['user_id'] = user_id
    if not slide_language is None:
        payload['slide_language'] = slide_language
    if not photo_permission is None:
        payload['photo_permission'] = photo_permission
    if not title is None:
        payload['title'] = title
    if not duration is None:
        payload['duration'] = duration
    if not spoken_language is None:
        payload['spoken_language'] = spoken_language
    if not memo is None:
        payload['memo'] = memo
    if not category is None:
        payload['category'] = category
    if not conference_id is None:
        payload['conference_id'] = conference_id
    if not material_level is None:
        payload['material_level'] = material_level
    if not slide_subtitles is None:
        payload['slide_subtitles'] = slide_subtitles
    uri = self.endpoint + "/session/create"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def lookup_session (self, id=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not id is None:
        payload['id'] = id
    uri = self.endpoint + "/session/lookup"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def delete_session (self, user_id=None, id=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not user_id is None:
        payload['user_id'] = user_id
    if not id is None:
        payload['id'] = id
    uri = self.endpoint + "/session/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def update_session (self, speaker_id=None, abstract=None, sort_order=None, tags=None, video_url=None, slide_url=None, status=None, title=None, user_id=None, has_interpretation=None, video_permission=None, photo_permission=None, slide_language=None, spoken_language=None, confirmed=None, duration=None, id=None, conference_id=None, slide_subtitles=None, material_level=None, category=None, memo=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not speaker_id is None:
        payload['speaker_id'] = speaker_id
    if not abstract is None:
        payload['abstract'] = abstract
    if not sort_order is None:
        payload['sort_order'] = sort_order
    if not tags is None:
        payload['tags'] = tags
    if not video_url is None:
        payload['video_url'] = video_url
    if not slide_url is None:
        payload['slide_url'] = slide_url
    if not status is None:
        payload['status'] = status
    if not title is None:
        payload['title'] = title
    if not user_id is None:
        payload['user_id'] = user_id
    if not has_interpretation is None:
        payload['has_interpretation'] = has_interpretation
    if not video_permission is None:
        payload['video_permission'] = video_permission
    if not photo_permission is None:
        payload['photo_permission'] = photo_permission
    if not slide_language is None:
        payload['slide_language'] = slide_language
    if not spoken_language is None:
        payload['spoken_language'] = spoken_language
    if not confirmed is None:
        payload['confirmed'] = confirmed
    if not duration is None:
        payload['duration'] = duration
    if not id is None:
        payload['id'] = id
    if not conference_id is None:
        payload['conference_id'] = conference_id
    if not slide_subtitles is None:
        payload['slide_subtitles'] = slide_subtitles
    if not material_level is None:
        payload['material_level'] = material_level
    if not category is None:
        payload['category'] = category
    if not memo is None:
        payload['memo'] = memo
    uri = self.endpoint + "/session/update"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def list_session_by_conference (self, conference_id=None, date=None):
    if conference_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not conference_id is None:
        payload['conference_id'] = conference_id
    if not date is None:
        payload['date'] = date
    uri = self.endpoint + "/schedule/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def create_question (self, body=None, user_id=None, session_id=None):
    if session_id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    if body is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not body is None:
        payload['body'] = body
    if not user_id is None:
        payload['user_id'] = user_id
    if not session_id is None:
        payload['session_id'] = session_id
    uri = self.endpoint + "/question/create"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def delete_question (self, id=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not id is None:
        payload['id'] = id
    uri = self.endpoint + "/question/delete"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def list_question (self, since=None, session_id=None, limit=None):
    if session_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not since is None:
        payload['since'] = since
    if not session_id is None:
        payload['session_id'] = session_id
    if not limit is None:
        payload['limit'] = limit
    uri = self.endpoint + "/question/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def create_session_survey_response (self, comment_good=None, overall_rating=None, speaker_knowledge=None, session_id=None, user_prior_knowledge=None, user_id=None, speaker_presentation=None, comment_improvement=None, material_quality=None):
    if session_id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    if user_prior_knowledge is None:
            raise "property \"" + required + "\" must be provided"
    if speaker_knowledge is None:
            raise "property \"" + required + "\" must be provided"
    if speaker_presentation is None:
            raise "property \"" + required + "\" must be provided"
    if material_quality is None:
            raise "property \"" + required + "\" must be provided"
    if overall_rating is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not comment_good is None:
        payload['comment_good'] = comment_good
    if not overall_rating is None:
        payload['overall_rating'] = overall_rating
    if not speaker_knowledge is None:
        payload['speaker_knowledge'] = speaker_knowledge
    if not session_id is None:
        payload['session_id'] = session_id
    if not user_prior_knowledge is None:
        payload['user_prior_knowledge'] = user_prior_knowledge
    if not user_id is None:
        payload['user_id'] = user_id
    if not speaker_presentation is None:
        payload['speaker_presentation'] = speaker_presentation
    if not comment_improvement is None:
        payload['comment_improvement'] = comment_improvement
    if not material_quality is None:
        payload['material_quality'] = material_quality
    uri = self.endpoint + "/survey_session_response/create"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

