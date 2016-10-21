"""OCTAV Client Library"""
"""DO NOT EDIT: This file was generated from ../spec/v1/api.json on Fri Oct 21 17:15:56 2016"""

import certifi
import json
import os
import re
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

  def extract_error(self, r):
    try:
      js = json.loads(r.data)
      if 'error' in js:
        self.error = js['error']
      elif 'message' in js:
        self.error = js['message']
    except BaseException as e:
      self.error = r.status

  def last_error(self):
    return self.error

  def last_response(self):
    return self.res

  def health_check (self):
    try:
        payload = {}
        hdrs = {}
        uri = '%s/' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def create_user (self, auth_user_id, auth_via, nickname, avatar_url=None, email=None, first_name=None, last_name=None, tshirt_size=None, **args):
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
        uri = '%s/v1/user/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def lookup_user (self, id):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v1/user/lookup' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def lookup_user_by_auth_user_id (self, auth_user_id, auth_via):
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
        uri = '%s/v1/user/lookup_user_by_auth_user_id' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def update_user (self, id, user_id, email=None, first_name=None, last_name=None, nickname=None, tshirt_size=None, **args):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
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
        patterns = [re.compile('first_name#[a-z]+'), re.compile('last_name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v1/user/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def delete_user (self, id, user_id):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if id is not None:
            payload['id'] = id
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/user/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def list_user (self, lang=None, limit=None, since=None):
    try:
        payload = {}
        hdrs = {}
        if lang is not None:
            payload['lang'] = lang
        if limit is not None:
            payload['limit'] = limit
        if since is not None:
            payload['since'] = since
        uri = '%s/v1/user/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def create_venue (self, address, name, user_id, latitude=None, longitude=None, **args):
    try:
        payload = {}
        hdrs = {}
        if address is None:
            raise MissingRequiredArgument('property address must be provided')
        payload['address'] = address
        if name is None:
            raise MissingRequiredArgument('property name must be provided')
        payload['name'] = name
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
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
        patterns = [re.compile('address#[a-z]+'), re.compile('name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v1/venue/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def list_venue (self, lang=None, limit=None, since=None):
    try:
        payload = {}
        hdrs = {}
        if lang is not None:
            payload['lang'] = lang
        if limit is not None:
            payload['limit'] = limit
        if since is not None:
            payload['since'] = since
        uri = '%s/v1/venue/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def lookup_venue (self, id):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v1/venue/lookup' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def update_venue (self, id, user_id):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if id is not None:
            payload['id'] = id
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/venue/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def delete_venue (self, id, user_id):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if id is not None:
            payload['id'] = id
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/venue/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def create_room (self, name, user_id, venue_id, capacity=None, **args):
    try:
        payload = {}
        hdrs = {}
        if name is None:
            raise MissingRequiredArgument('property name must be provided')
        payload['name'] = name
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if venue_id is None:
            raise MissingRequiredArgument('property venue_id must be provided')
        payload['venue_id'] = venue_id
        if capacity is not None:
            payload['capacity'] = capacity
        if name is not None:
            payload['name'] = name
        if user_id is not None:
            payload['user_id'] = user_id
        if venue_id is not None:
            payload['venue_id'] = venue_id
        patterns = [re.compile('name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v1/room/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def update_room (self, id, user_id, capacity=None, name=None, venue_id=None, **args):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
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
        patterns = [re.compile('name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v1/room/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def lookup_room (self, id):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v1/room/lookup' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def delete_room (self, id, user_id):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if id is not None:
            payload['id'] = id
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/room/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def list_room (self, venue_id, lang=None, limit=None):
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
        uri = '%s/v1/room/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def create_conference_series (self, slug, title, user_id, description=None):
    try:
        payload = {}
        hdrs = {}
        if slug is None:
            raise MissingRequiredArgument('property slug must be provided')
        payload['slug'] = slug
        if title is None:
            raise MissingRequiredArgument('property title must be provided')
        payload['title'] = title
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if description is not None:
            payload['description'] = description
        if slug is not None:
            payload['slug'] = slug
        if title is not None:
            payload['title'] = title
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/conference_series/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def lookup_conference_series (self, id, lang=None):
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
        uri = '%s/v1/conference_series/lookup' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def list_conference_series (self, limit=None, since=None):
    try:
        payload = {}
        hdrs = {}
        if limit is not None:
            payload['limit'] = limit
        if since is not None:
            payload['since'] = since
        uri = '%s/v1/conference_series/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def add_conference_series_admin (self, admin_id, series_id, user_id):
    try:
        payload = {}
        hdrs = {}
        if admin_id is None:
            raise MissingRequiredArgument('property admin_id must be provided')
        payload['admin_id'] = admin_id
        if series_id is None:
            raise MissingRequiredArgument('property series_id must be provided')
        payload['series_id'] = series_id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if admin_id is not None:
            payload['admin_id'] = admin_id
        if series_id is not None:
            payload['series_id'] = series_id
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/conference_series/admin/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def create_conference (self, series_id, slug, title, user_id, cfp_lead_text=None, cfp_post_submit_instructions=None, cfp_pre_submit_instructions=None, description=None, sub_title=None, timezone=None):
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
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if cfp_lead_text is not None:
            payload['cfp_lead_text'] = cfp_lead_text
        if cfp_post_submit_instructions is not None:
            payload['cfp_post_submit_instructions'] = cfp_post_submit_instructions
        if cfp_pre_submit_instructions is not None:
            payload['cfp_pre_submit_instructions'] = cfp_pre_submit_instructions
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
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/conference/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def get_conference_schedule (self, conference_id):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        uri = '%s/v1/conference/schedule.ics' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def add_conference_credential (self, conference_id, data, type, user_id):
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
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if data is not None:
            payload['data'] = data
        if type is not None:
            payload['type'] = type
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/conference/credentials/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def tweet_as_conference (self, conference_id, tweet, user_id):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if tweet is None:
            raise MissingRequiredArgument('property tweet must be provided')
        payload['tweet'] = tweet
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if tweet is not None:
            payload['tweet'] = tweet
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/conference/tweet' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def add_conference_date (self, conference_id, date, user_id):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if date is None:
            raise MissingRequiredArgument('property date must be provided')
        payload['date'] = date
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if date is not None:
            payload['date'] = date
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/conference/date/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def delete_conference_date (self, conference_id, date, user_id):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if date is None:
            raise MissingRequiredArgument('property date must be provided')
        payload['date'] = date
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if date is not None:
            payload['date'] = date
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/conference/date/delete' % self.endpoint
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def add_conference_admin (self, admin_id, conference_id, user_id):
    try:
        payload = {}
        hdrs = {}
        if admin_id is None:
            raise MissingRequiredArgument('property admin_id must be provided')
        payload['admin_id'] = admin_id
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if admin_id is not None:
            payload['admin_id'] = admin_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/conference/admin/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def delete_conference_admin (self, admin_id, conference_id, user_id):
    try:
        payload = {}
        hdrs = {}
        if admin_id is None:
            raise MissingRequiredArgument('property admin_id must be provided')
        payload['admin_id'] = admin_id
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if admin_id is not None:
            payload['admin_id'] = admin_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/conference/admin/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def add_conference_venue (self, conference_id, user_id, venue_id):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if venue_id is None:
            raise MissingRequiredArgument('property venue_id must be provided')
        payload['venue_id'] = venue_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if user_id is not None:
            payload['user_id'] = user_id
        if venue_id is not None:
            payload['venue_id'] = venue_id
        uri = '%s/v1/conference/venue/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def delete_conference_venue (self, conference_id, user_id, venue_id):
    try:
        payload = {}
        hdrs = {}
        if conference_id is None:
            raise MissingRequiredArgument('property conference_id must be provided')
        payload['conference_id'] = conference_id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if venue_id is None:
            raise MissingRequiredArgument('property venue_id must be provided')
        payload['venue_id'] = venue_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if user_id is not None:
            payload['user_id'] = user_id
        if venue_id is not None:
            payload['venue_id'] = venue_id
        uri = '%s/v1/conference/venue/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def add_session_type (self, abstract, conference_id, duration, name, user_id, submission_end=None, submission_start=None, **args):
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
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
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
        if user_id is not None:
            payload['user_id'] = user_id
        patterns = [re.compile('abstract#[a-z]+'), re.compile('name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v1/conference/session_type/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def delete_session_type (self, id, user_id):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if id is not None:
            payload['id'] = id
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/session_type/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def lookup_session_type (self, id, lang=None):
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
        uri = '%s/v1/session_type/lookup' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def update_session_type (self, id, user_id, abstract=None, duration=None, name=None, submission_end=None, submission_start=None, **args):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if abstract is not None:
            payload['abstract'] = abstract
        if duration is not None:
            payload['duration'] = duration
        if id is not None:
            payload['id'] = id
        if name is not None:
            payload['name'] = name
        if submission_end is not None:
            payload['submission_end'] = submission_end
        if submission_start is not None:
            payload['submission_start'] = submission_start
        if user_id is not None:
            payload['user_id'] = user_id
        patterns = [re.compile('abstract#[a-z]+'), re.compile('name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v1/session_type/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def list_session_types_by_conference (self, conference_id=None, lang=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if lang is not None:
            payload['lang'] = lang
        uri = '%s/v1/session_type/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def lookup_conference (self, id, lang=None):
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
        uri = '%s/v1/conference/lookup' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def lookup_conference_by_slug (self, slug, lang=None):
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
        uri = '%s/v1/conference/lookup_by_slug' % self.endpoint
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def list_conferences_by_organizer (self, lang=None, limit=None, organizer_id=None, since=None):
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
        uri = '%s/v1/conference/list_by_organizer' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def list_conference (self, lang=None, limit=None, range_end=None, range_start=None, since=None, status=None):
    try:
        payload = {}
        hdrs = {}
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
        uri = '%s/v1/conference/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def update_conference (self, id, user_id, cfp_lead_text=None, cfp_post_submit_instructions=None, cfp_pre_submit_instructions=None, description=None, slug=None, status=None, sub_title=None, timezone=None, title=None, **args):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if cfp_lead_text is not None:
            payload['cfp_lead_text'] = cfp_lead_text
        if cfp_post_submit_instructions is not None:
            payload['cfp_post_submit_instructions'] = cfp_post_submit_instructions
        if cfp_pre_submit_instructions is not None:
            payload['cfp_pre_submit_instructions'] = cfp_pre_submit_instructions
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
        if timezone is not None:
            payload['timezone'] = timezone
        if title is not None:
            payload['title'] = title
        if user_id is not None:
            payload['user_id'] = user_id
        patterns = [re.compile('cfp_lead_text#[a-z]+'), re.compile('cfp_post_submit_instructions#[a-z]+'), re.compile('cfp_pre_submit_instructions#[a-z]+'), re.compile('description#[a-z]+'), re.compile('title#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v1/conference/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def delete_conference_series (self, id, user_id):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if id is not None:
            payload['id'] = id
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/conference_series/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def delete_conference (self, id):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v1/conference/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def create_session (self, conference_id, session_type_id, speaker_id, user_id, abstract=None, category=None, material_level=None, materials_release=None, memo=None, photo_release=None, recording_release=None, slide_language=None, slide_subtitles=None, slide_url=None, spoken_language=None, tags=None, title=None, video_url=None, **args):
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
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
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
        if user_id is not None:
            payload['user_id'] = user_id
        if video_url is not None:
            payload['video_url'] = video_url
        patterns = [re.compile('abstract#[a-z]+'), re.compile('title#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v1/session/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def lookup_session (self, id, lang=None):
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
        uri = '%s/v1/session/lookup' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def delete_session (self, id, user_id):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if id is not None:
            payload['id'] = id
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/session/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def update_session (self, id, user_id, abstract=None, category=None, conference_id=None, confirmed=None, duration=None, has_interpretation=None, material_level=None, materials_release=None, memo=None, photo_release=None, recording_release=None, session_type_id=None, slide_language=None, slide_subtitles=None, slide_url=None, sort_order=None, speaker_id=None, spoken_language=None, starts_on=None, status=None, tags=None, title=None, video_url=None, **args):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
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
        if user_id is not None:
            payload['user_id'] = user_id
        if video_url is not None:
            payload['video_url'] = video_url
        patterns = [re.compile('abstract#[a-z]+'), re.compile('title#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v1/session/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def list_sessions (self, conference_id=None, date=None, lang=None, limit=None, since=None, speaker_id=None, status=None):
    try:
        payload = {}
        hdrs = {}
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if date is not None:
            payload['date'] = date
        if lang is not None:
            payload['lang'] = lang
        if limit is not None:
            payload['limit'] = limit
        if since is not None:
            payload['since'] = since
        if speaker_id is not None:
            payload['speaker_id'] = speaker_id
        if status is not None:
            payload['status'] = status
        uri = '%s/v1/session/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def create_question (self, body, session_id, user_id):
    try:
        payload = {}
        hdrs = {}
        if body is None:
            raise MissingRequiredArgument('property body must be provided')
        payload['body'] = body
        if session_id is None:
            raise MissingRequiredArgument('property session_id must be provided')
        payload['session_id'] = session_id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if body is not None:
            payload['body'] = body
        if session_id is not None:
            payload['session_id'] = session_id
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/question/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def delete_question (self, id):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if id is not None:
            payload['id'] = id
        uri = '%s/v1/question/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def list_question (self, session_id, limit=None, since=None):
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
        uri = '%s/v1/question/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def create_session_survey_response (self, material_quality, overall_rating, session_id, speaker_knowledge, speaker_presentation, user_id, user_prior_knowledge, comment_good=None, comment_improvement=None):
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
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
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
        if user_id is not None:
            payload['user_id'] = user_id
        if user_prior_knowledge is not None:
            payload['user_prior_knowledge'] = user_prior_knowledge
        uri = '%s/v1/survey_session_response/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def add_featured_speaker (self, conference_id, description, display_name, avatar_url=None, speaker_id=None, user_id=None, **args):
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
        if user_id is not None:
            payload['user_id'] = user_id
        patterns = [re.compile('description#[a-z]+'), re.compile('display_name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v1/featured_speaker/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def lookup_featured_speaker (self, id, lang=None):
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
        uri = '%s/v1/featured_speaker/lookup' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def list_featured_speakers (self, conference_id=None, lang=None, limit=None, since=None):
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
        uri = '%s/v1/featured_speaker/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def update_featured_speaker (self, id, user_id, avatar_url=None, description=None, display_name=None, speaker_id=None, **args):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
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
        patterns = [re.compile('description#[a-z]+'), re.compile('display_name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v1/featured_speaker/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def delete_featured_speaker (self, id, user_id):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if id is not None:
            payload['id'] = id
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/featured_speaker/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def add_sponsor (self, conference_id, group_name, name, url, user_id, sort_order=None, **args):
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
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if conference_id is not None:
            payload['conference_id'] = conference_id
        if group_name is not None:
            payload['group_name'] = group_name
        if name is not None:
            payload['name'] = name
        if sort_order is not None:
            payload['sort_order'] = sort_order
        if url is not None:
            payload['url'] = url
        if user_id is not None:
            payload['user_id'] = user_id
        patterns = [re.compile('name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v1/sponsor/add' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def lookup_sponsor (self, id, lang=None):
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
        uri = '%s/v1/sponsor/lookup' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def list_sponsors (self, conference_id=None, lang=None, limit=None, since=None):
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
        uri = '%s/v1/sponsor/list' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        qs = urlencode(payload, True)
        if self.debug:
            print('GET %s?%s' % (uri, qs))
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

  def update_sponsor (self, id, user_id, group_name=None, name=None, sort_order=None, url=None, **args):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if group_name is not None:
            payload['group_name'] = group_name
        if id is not None:
            payload['id'] = id
        if name is not None:
            payload['name'] = name
        if sort_order is not None:
            payload['sort_order'] = sort_order
        if url is not None:
            payload['url'] = url
        if user_id is not None:
            payload['user_id'] = user_id
        patterns = [re.compile('name#[a-z]+')]
        for key in args:
            for p in patterns:
                if p.match(key):
                    payload[key] = args[key]
        uri = '%s/v1/sponsor/update' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def delete_sponsor (self, id, user_id):
    try:
        payload = {}
        hdrs = {}
        if id is None:
            raise MissingRequiredArgument('property id must be provided')
        payload['id'] = id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if id is not None:
            payload['id'] = id
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/sponsor/delete' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def create_temporary_email (self, email, target_id, user_id):
    try:
        payload = {}
        hdrs = {}
        if email is None:
            raise MissingRequiredArgument('property email must be provided')
        payload['email'] = email
        if target_id is None:
            raise MissingRequiredArgument('property target_id must be provided')
        payload['target_id'] = target_id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if email is not None:
            payload['email'] = email
        if target_id is not None:
            payload['target_id'] = target_id
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/email/create' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

  def confirm_temporary_email (self, confirmation_key, target_id, user_id):
    try:
        payload = {}
        hdrs = {}
        if confirmation_key is None:
            raise MissingRequiredArgument('property confirmation_key must be provided')
        payload['confirmation_key'] = confirmation_key
        if target_id is None:
            raise MissingRequiredArgument('property target_id must be provided')
        payload['target_id'] = target_id
        if user_id is None:
            raise MissingRequiredArgument('property user_id must be provided')
        payload['user_id'] = user_id
        if confirmation_key is not None:
            payload['confirmation_key'] = confirmation_key
        if target_id is not None:
            payload['target_id'] = target_id
        if user_id is not None:
            payload['user_id'] = user_id
        uri = '%s/v1/email/confirm' % self.endpoint
        hdrs = urllib3.util.make_headers(
            basic_auth='%s:%s' % (self.key, self.secret),
        )
        if self.debug:
            print('POST %s' % uri)
        hdrs['Content-Type']= 'application/json'
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

