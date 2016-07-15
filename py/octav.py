"""OCTAV Client Library"""
"""DO NOT EDIT: This file was generated from ../spec/v1/api.json on Fri Jul 15 10:34:03 2016"""

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
    serlf.error = None
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

  def create_user (self, avatar_url=None, tshirt_size=None, nickname=None, last_name=None, auth_user_id=None, auth_via=None, email=None, first_name=None):
    if nickname is None:
            raise "property \"" + required + "\" must be provided"
    if auth_via is None:
            raise "property \"" + required + "\" must be provided"
    if auth_user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not avatar_url is None:
        payload['avatar_url'] = avatar_url
    if not tshirt_size is None:
        payload['tshirt_size'] = tshirt_size
    if not nickname is None:
        payload['nickname'] = nickname
    if not last_name is None:
        payload['last_name'] = last_name
    if not auth_user_id is None:
        payload['auth_user_id'] = auth_user_id
    if not auth_via is None:
        payload['auth_via'] = auth_via
    if not email is None:
        payload['email'] = email
    if not first_name is None:
        payload['first_name'] = first_name
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

  def lookup_user_by_auth_user_id (self, auth_user_id=None, auth_via=None):
    if auth_via is None:
            raise "property \"" + required + "\" must be provided"
    if auth_user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not auth_user_id is None:
        payload['auth_user_id'] = auth_user_id
    if not auth_via is None:
        payload['auth_via'] = auth_via
    uri = self.endpoint + "/user/lookup_by_auth_user_id"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def update_user (self, email=None, user_id=None, tshirt_size=None, last_name=None, id=None, nickname=None, first_name=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not email is None:
        payload['email'] = email
    if not user_id is None:
        payload['user_id'] = user_id
    if not tshirt_size is None:
        payload['tshirt_size'] = tshirt_size
    if not last_name is None:
        payload['last_name'] = last_name
    if not id is None:
        payload['id'] = id
    if not nickname is None:
        payload['nickname'] = nickname
    if not first_name is None:
        payload['first_name'] = first_name
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

  def list_user (self, lang=None, since=None, limit=None):
    payload = {}
    if not lang is None:
        payload['lang'] = lang
    if not since is None:
        payload['since'] = since
    if not limit is None:
        payload['limit'] = limit
    uri = self.endpoint + "/user/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def create_venue (self, latitude=None, user_id=None, address=None, name=None, longitude=None):
    if name is None:
            raise "property \"" + required + "\" must be provided"
    if address is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not latitude is None:
        payload['latitude'] = latitude
    if not user_id is None:
        payload['user_id'] = user_id
    if not address is None:
        payload['address'] = address
    if not name is None:
        payload['name'] = name
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

  def list_venue (self, since=None, limit=None, lang=None):
    payload = {}
    if not since is None:
        payload['since'] = since
    if not limit is None:
        payload['limit'] = limit
    if not lang is None:
        payload['lang'] = lang
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

  def create_room (self, venue_id=None, user_id=None, name=None, capacity=None):
    if venue_id is None:
            raise "property \"" + required + "\" must be provided"
    if name is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not venue_id is None:
        payload['venue_id'] = venue_id
    if not user_id is None:
        payload['user_id'] = user_id
    if not name is None:
        payload['name'] = name
    if not capacity is None:
        payload['capacity'] = capacity
    uri = self.endpoint + "/room/create"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def update_room (self, name=None, capacity=None, user_id=None, id=None, venue_id=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not name is None:
        payload['name'] = name
    if not capacity is None:
        payload['capacity'] = capacity
    if not user_id is None:
        payload['user_id'] = user_id
    if not id is None:
        payload['id'] = id
    if not venue_id is None:
        payload['venue_id'] = venue_id
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

  def list_room (self, venue_id=None, limit=None, lang=None):
    if venue_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not venue_id is None:
        payload['venue_id'] = venue_id
    if not limit is None:
        payload['limit'] = limit
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

  def list_conference_series (self, since=None, limit=None):
    payload = {}
    if not since is None:
        payload['since'] = since
    if not limit is None:
        payload['limit'] = limit
    uri = self.endpoint + "/conference_series/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def add_conference_series_admin (self, series_id=None, user_id=None, admin_id=None):
    if series_id is None:
            raise "property \"" + required + "\" must be provided"
    if admin_id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not series_id is None:
        payload['series_id'] = series_id
    if not user_id is None:
        payload['user_id'] = user_id
    if not admin_id is None:
        payload['admin_id'] = admin_id
    uri = self.endpoint + "/conference_series/admin/add"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def create_conference (self, series_id=None, user_id=None, description=None, slug=None, title=None, sub_title=None):
    if series_id is None:
            raise "property \"" + required + "\" must be provided"
    if title is None:
            raise "property \"" + required + "\" must be provided"
    if slug is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not series_id is None:
        payload['series_id'] = series_id
    if not user_id is None:
        payload['user_id'] = user_id
    if not description is None:
        payload['description'] = description
    if not slug is None:
        payload['slug'] = slug
    if not title is None:
        payload['title'] = title
    if not sub_title is None:
        payload['sub_title'] = sub_title
    uri = self.endpoint + "/conference/create"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def add_conference_dates (self, conference_id=None, dates=None, user_id=None):
    if conference_id is None:
            raise "property \"" + required + "\" must be provided"
    if dates is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not conference_id is None:
        payload['conference_id'] = conference_id
    if not dates is None:
        payload['dates'] = dates
    if not user_id is None:
        payload['user_id'] = user_id
    uri = self.endpoint + "/conference/dates/add"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_conference_dates (self, conference_id=None, dates=None, user_id=None):
    if conference_id is None:
            raise "property \"" + required + "\" must be provided"
    if dates is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not conference_id is None:
        payload['conference_id'] = conference_id
    if not dates is None:
        payload['dates'] = dates
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

  def add_conference_admin (self, user_id=None, admin_id=None, conference_id=None):
    if conference_id is None:
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
    if not conference_id is None:
        payload['conference_id'] = conference_id
    uri = self.endpoint + "/conference/admin/add"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_conference_admin (self, conference_id=None, user_id=None, admin_id=None):
    if conference_id is None:
            raise "property \"" + required + "\" must be provided"
    if admin_id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not conference_id is None:
        payload['conference_id'] = conference_id
    if not user_id is None:
        payload['user_id'] = user_id
    if not admin_id is None:
        payload['admin_id'] = admin_id
    uri = self.endpoint + "/conference/admin/delete"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def add_conference_venue (self, venue_id=None, conference_id=None, user_id=None):
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

  def lookup_conference_by_slug (self, lang=None, slug=None):
    if slug is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not lang is None:
        payload['lang'] = lang
    if not slug is None:
        payload['slug'] = slug
    uri = self.endpoint + "/conference/lookup_by_slug"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def list_conference (self, since=None, limit=None, range_end=None, lang=None, range_start=None):
    payload = {}
    if not since is None:
        payload['since'] = since
    if not limit is None:
        payload['limit'] = limit
    if not range_end is None:
        payload['range_end'] = range_end
    if not lang is None:
        payload['lang'] = lang
    if not range_start is None:
        payload['range_start'] = range_start
    uri = self.endpoint + "/conference/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def update_conference (self, slug=None, user_id=None, description=None, starts_on=None, id=None, sub_title=None, title=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not slug is None:
        payload['slug'] = slug
    if not user_id is None:
        payload['user_id'] = user_id
    if not description is None:
        payload['description'] = description
    if not starts_on is None:
        payload['starts_on'] = starts_on
    if not id is None:
        payload['id'] = id
    if not sub_title is None:
        payload['sub_title'] = sub_title
    if not title is None:
        payload['title'] = title
    uri = self.endpoint + "/conference/update"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def delete_conference_series (self, id=None, user_id=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not id is None:
        payload['id'] = id
    if not user_id is None:
        payload['user_id'] = user_id
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

  def create_session (self, slide_language=None, spoken_language=None, slide_url=None, title=None, tags=None, duration=None, abstract=None, material_level=None, video_permission=None, speaker_id=None, slide_subtitles=None, memo=None, conference_id=None, photo_permission=None, category=None, user_id=None, video_url=None):
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
    if not slide_language is None:
        payload['slide_language'] = slide_language
    if not spoken_language is None:
        payload['spoken_language'] = spoken_language
    if not slide_url is None:
        payload['slide_url'] = slide_url
    if not title is None:
        payload['title'] = title
    if not tags is None:
        payload['tags'] = tags
    if not duration is None:
        payload['duration'] = duration
    if not abstract is None:
        payload['abstract'] = abstract
    if not material_level is None:
        payload['material_level'] = material_level
    if not video_permission is None:
        payload['video_permission'] = video_permission
    if not speaker_id is None:
        payload['speaker_id'] = speaker_id
    if not slide_subtitles is None:
        payload['slide_subtitles'] = slide_subtitles
    if not memo is None:
        payload['memo'] = memo
    if not conference_id is None:
        payload['conference_id'] = conference_id
    if not photo_permission is None:
        payload['photo_permission'] = photo_permission
    if not category is None:
        payload['category'] = category
    if not user_id is None:
        payload['user_id'] = user_id
    if not video_url is None:
        payload['video_url'] = video_url
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

  def update_session (self, video_url=None, user_id=None, category=None, photo_permission=None, has_interpretation=None, memo=None, conference_id=None, slide_subtitles=None, id=None, speaker_id=None, video_permission=None, sort_order=None, material_level=None, abstract=None, tags=None, duration=None, title=None, confirmed=None, slide_url=None, status=None, spoken_language=None, slide_language=None):
    if id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not video_url is None:
        payload['video_url'] = video_url
    if not user_id is None:
        payload['user_id'] = user_id
    if not category is None:
        payload['category'] = category
    if not photo_permission is None:
        payload['photo_permission'] = photo_permission
    if not has_interpretation is None:
        payload['has_interpretation'] = has_interpretation
    if not memo is None:
        payload['memo'] = memo
    if not conference_id is None:
        payload['conference_id'] = conference_id
    if not slide_subtitles is None:
        payload['slide_subtitles'] = slide_subtitles
    if not id is None:
        payload['id'] = id
    if not speaker_id is None:
        payload['speaker_id'] = speaker_id
    if not video_permission is None:
        payload['video_permission'] = video_permission
    if not sort_order is None:
        payload['sort_order'] = sort_order
    if not material_level is None:
        payload['material_level'] = material_level
    if not abstract is None:
        payload['abstract'] = abstract
    if not tags is None:
        payload['tags'] = tags
    if not duration is None:
        payload['duration'] = duration
    if not title is None:
        payload['title'] = title
    if not confirmed is None:
        payload['confirmed'] = confirmed
    if not slide_url is None:
        payload['slide_url'] = slide_url
    if not status is None:
        payload['status'] = status
    if not spoken_language is None:
        payload['spoken_language'] = spoken_language
    if not slide_language is None:
        payload['slide_language'] = slide_language
    uri = self.endpoint + "/session/update"
    if self.debug:
        print("POST " + uri)
    res = self.session.post(uri, auth=(self.key, self.secret), json=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

  def list_session_by_conference (self, date=None, conference_id=None):
    if conference_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not date is None:
        payload['date'] = date
    if not conference_id is None:
        payload['conference_id'] = conference_id
    uri = self.endpoint + "/schedule/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def create_question (self, user_id=None, session_id=None, body=None):
    if session_id is None:
            raise "property \"" + required + "\" must be provided"
    if user_id is None:
            raise "property \"" + required + "\" must be provided"
    if body is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not user_id is None:
        payload['user_id'] = user_id
    if not session_id is None:
        payload['session_id'] = session_id
    if not body is None:
        payload['body'] = body
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

  def list_question (self, session_id=None, limit=None, since=None):
    if session_id is None:
            raise "property \"" + required + "\" must be provided"
    payload = {}
    if not session_id is None:
        payload['session_id'] = session_id
    if not limit is None:
        payload['limit'] = limit
    if not since is None:
        payload['since'] = since
    uri = self.endpoint + "/question/list"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return res.json()

  def create_session_survey_response (self, overall_rating=None, comment_improvement=None, comment_good=None, material_quality=None, user_prior_knowledge=None, speaker_knowledge=None, user_id=None, session_id=None, speaker_presentation=None):
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
    if not overall_rating is None:
        payload['overall_rating'] = overall_rating
    if not comment_improvement is None:
        payload['comment_improvement'] = comment_improvement
    if not comment_good is None:
        payload['comment_good'] = comment_good
    if not material_quality is None:
        payload['material_quality'] = material_quality
    if not user_prior_knowledge is None:
        payload['user_prior_knowledge'] = user_prior_knowledge
    if not speaker_knowledge is None:
        payload['speaker_knowledge'] = speaker_knowledge
    if not user_id is None:
        payload['user_id'] = user_id
    if not session_id is None:
        payload['session_id'] = session_id
    if not speaker_presentation is None:
        payload['speaker_presentation'] = speaker_presentation
    uri = self.endpoint + "/survey_session_response/create"
    if self.debug:
        print("GET " + uri)
    res = self.session.get(uri, params=payload)
    if res.status_code != 200:
        self.extract_error(res)
        return None
    return True

