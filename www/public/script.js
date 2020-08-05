/*global PushStream WebSocketWrapper EventSourceWrapper EventSource*/
/*jshint evil: true, plusplus: false, regexp: false */
/**
The MIT License (MIT)

Copyright (c) 2010-2014 Wandenberg Peixoto <wandenberg@gmail.com>, Rogério Carvalho Schneider <stockrt@gmail.com>

This file is part of Nginx Push Stream Module.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

pushstream.js

Created: Nov 01, 2011
Authors: Wandenberg Peixoto <wandenberg@gmail.com>, Rogério Carvalho Schneider <stockrt@gmail.com>
 */
(function (window, document, undefined) {
    "use strict";
  
    /* prevent duplicate declaration */
    if (window.PushStream) { return; }
  
    var Utils = {};
  
    var days = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
    var months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
  
    var valueToTwoDigits = function (value) {
      return ((value < 10) ? '0' : '') + value;
    };
  
    Utils.dateToUTCString = function (date) {
      var time = valueToTwoDigits(date.getUTCHours()) + ':' + valueToTwoDigits(date.getUTCMinutes()) + ':' + valueToTwoDigits(date.getUTCSeconds());
      return days[date.getUTCDay()] + ', ' + valueToTwoDigits(date.getUTCDate()) + ' ' + months[date.getUTCMonth()] + ' ' + date.getUTCFullYear() + ' ' + time + ' GMT';
    };
  
    var extend = function () {
      var object = arguments[0] || {};
      for (var i = 0; i < arguments.length; i++) {
        var settings = arguments[i];
        for (var attr in settings) {
          if (!settings.hasOwnProperty || settings.hasOwnProperty(attr)) {
            object[attr] = settings[attr];
          }
        }
      }
      return object;
    };
  
    var validChars  = /^[\],:{}\s]*$/,
        validEscape = /\\(?:["\\\/bfnrt]|u[0-9a-fA-F]{4})/g,
        validTokens = /"[^"\\\n\r]*"|true|false|null|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?/g,
        validBraces = /(?:^|:|,)(?:\s*\[)+/g;
  
    var trim = function(value) {
      return value.replace(/^\s*/, "").replace(/\s*$/, "");
    };
  
    Utils.parseJSON = function(data) {
      if (!data || !isString(data)) {
        return null;
      }
  
      // Make sure leading/trailing whitespace is removed (IE can't handle it)
      data = trim(data);
  
      // Attempt to parse using the native JSON parser first
      if (window.JSON && window.JSON.parse) {
        try {
          return window.JSON.parse( data );
        } catch(e) {
          throw "Invalid JSON: " + data;
        }
      }
  
      // Make sure the incoming data is actual JSON
      // Logic borrowed from http://json.org/json2.js
      if (validChars.test(data.replace(validEscape, "@").replace( validTokens, "]").replace( validBraces, "")) ) {
        return (new Function("return " + data))();
      }
  
      throw "Invalid JSON: " + data;
    };
  
    var getControlParams = function(pushstream) {
      var data = {};
      data[pushstream.tagArgument] = "";
      data[pushstream.timeArgument] = "";
      data[pushstream.eventIdArgument] = "";
      if (pushstream.messagesControlByArgument) {
        data[pushstream.tagArgument] = Number(pushstream._etag);
        if (pushstream._lastModified) {
          data[pushstream.timeArgument] = pushstream._lastModified;
        } else if (pushstream._lastEventId) {
          data[pushstream.eventIdArgument] = pushstream._lastEventId;
        }
      }
      return data;
    };
  
    var getTime = function() {
      return (new Date()).getTime();
    };
  
    var currentTimestampParam = function() {
      return { "_" : getTime() };
    };
  
    var objectToUrlParams = function(settings) {
      var params = settings;
      if (typeof(settings) === 'object') {
        params = '';
        for (var attr in settings) {
          if (!settings.hasOwnProperty || settings.hasOwnProperty(attr)) {
            params += '&' + attr + '=' + escapeText(settings[attr]);
          }
        }
        params = params.substring(1);
      }
  
      return params || '';
    };
  
    var addParamsToUrl = function(url, params) {
      return url + ((url.indexOf('?') < 0) ? '?' : '&') + objectToUrlParams(params);
    };
  
    var isArray = Array.isArray || function(obj) {
      return Object.prototype.toString.call(obj) === '[object Array]';
    };
  
    var isString = function(obj) {
      return Object.prototype.toString.call(obj) === '[object String]';
    };
  
    var isDate = function(obj) {
      return Object.prototype.toString.call(obj) === '[object Date]';
    };
  
    var Log4js = {
      logger: null,
      debug : function() { if  (PushStream.LOG_LEVEL === 'debug')                                         { Log4js._log.apply(Log4js._log, arguments); }},
      info  : function() { if ((PushStream.LOG_LEVEL === 'info')  || (PushStream.LOG_LEVEL === 'debug'))  { Log4js._log.apply(Log4js._log, arguments); }},
      error : function() {                                                                                  Log4js._log.apply(Log4js._log, arguments); },
      _initLogger : function() {
        var console = window.console;
        if (console && console.log) {
          if (console.log.apply) {
            Log4js.logger = console.log;
          } else if ((typeof console.log === "object") && Function.prototype.bind) {
            Log4js.logger = Function.prototype.bind.call(console.log, console);
          } else if ((typeof console.log === "object") && Function.prototype.call) {
            Log4js.logger = function() {
              Function.prototype.call.call(console.log, console, Array.prototype.slice.call(arguments));
            };
          }
        }
      },
      _log  : function() {
        if (!Log4js.logger) {
          Log4js._initLogger();
        }
  
        if (Log4js.logger) {
          try {
            Log4js.logger.apply(window.console, arguments);
          } catch(e) {
            Log4js._initLogger();
            if (Log4js.logger) {
              Log4js.logger.apply(window.console, arguments);
            }
          }
        }
  
        var logElement = document.getElementById(PushStream.LOG_OUTPUT_ELEMENT_ID);
        if (logElement) {
          var str = '';
          for (var i = 0; i < arguments.length; i++) {
            str += arguments[i] + " ";
          }
          logElement.innerHTML += str + '\n';
  
          var lines = logElement.innerHTML.split('\n');
          if (lines.length > 100) {
            logElement.innerHTML = lines.slice(-100).join('\n');
          }
        }
      }
    };
  
    var Ajax = {
      _getXHRObject : function(crossDomain) {
        var xhr = false;
        if (crossDomain) {
          try { xhr = new window.XDomainRequest(); } catch (e) { }
          if (xhr) {
            return xhr;
          }
        }
  
        try { xhr = new window.XMLHttpRequest(); }
        catch (e1) {
          try { xhr = new window.ActiveXObject("Msxml2.XMLHTTP"); }
          catch (e2) {
            try { xhr = new window.ActiveXObject("Microsoft.XMLHTTP"); }
            catch (e3) {
              xhr = false;
            }
          }
        }
        return xhr;
      },
  
      _send : function(settings, post) {
        settings = settings || {};
        settings.timeout = settings.timeout || 30000;
        var xhr = Ajax._getXHRObject(settings.crossDomain);
        if (!xhr||!settings.url) { return; }
  
        Ajax.clear(settings);
  
        settings.xhr = xhr;
  
        if (window.XDomainRequest && (xhr instanceof window.XDomainRequest)) {
          xhr.onload = function () {
            if (settings.afterReceive) { settings.afterReceive(xhr); }
            if (settings.success) { settings.success(xhr.responseText); }
          };
  
          xhr.onerror = xhr.ontimeout = function () {
            if (settings.afterReceive) { settings.afterReceive(xhr); }
            if (settings.error) { settings.error(xhr.status); }
          };
        } else {
          xhr.onreadystatechange = function () {
            if (xhr.readyState === 4) {
              Ajax.clear(settings);
              if (settings.afterReceive) { settings.afterReceive(xhr); }
              if(xhr.status === 200) {
                if (settings.success) { settings.success(xhr.responseText); }
              } else {
                if (settings.error) { settings.error(xhr.status); }
              }
            }
          };
        }
  
        if (settings.beforeOpen) { settings.beforeOpen(xhr); }
  
        var params = {};
        var body = null;
        var method = "GET";
        if (post) {
          body = objectToUrlParams(settings.data);
          method = "POST";
        } else {
          params = settings.data || {};
        }
  
        xhr.open(method, addParamsToUrl(settings.url, extend({}, params, currentTimestampParam())), true);
  
        if (settings.beforeSend) { settings.beforeSend(xhr); }
  
        var onerror = function() {
          Ajax.clear(settings);
          try { xhr.abort(); } catch (e) { /* ignore error on closing */ }
          settings.error(304);
        };
  
        if (post) {
          if (xhr.setRequestHeader) {
            xhr.setRequestHeader("Accept", "application/json");
            xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
          }
        } else {
          settings.timeoutId = window.setTimeout(onerror, settings.timeout + 2000);
        }
  
        xhr.send(body);
        return xhr;
      },
  
      _clear_xhr : function(xhr) {
        if (xhr) {
          xhr.onreadystatechange = null;
        }
      },
  
      _clear_script : function(script) {
        // Handling memory leak in IE, removing and dereference the script
        if (script) {
          script.onerror = script.onload = script.onreadystatechange = null;
          if (script.parentNode) { script.parentNode.removeChild(script); }
        }
      },
  
      _clear_timeout : function(settings) {
        settings.timeoutId = clearTimer(settings.timeoutId);
      },
  
      _clear_jsonp : function(settings) {
        var callbackFunctionName = settings.data.callback;
        if (callbackFunctionName) {
          window[callbackFunctionName] = function() { window[callbackFunctionName] = null; };
        }
      },
  
      clear : function(settings) {
        Ajax._clear_timeout(settings);
        Ajax._clear_jsonp(settings);
        Ajax._clear_script(document.getElementById(settings.scriptId));
        Ajax._clear_xhr(settings.xhr);
      },
  
      jsonp : function(settings) {
        settings.timeout = settings.timeout || 30000;
        Ajax.clear(settings);
  
        var head = document.head || document.getElementsByTagName("head")[0];
        var script = document.createElement("script");
        var startTime = getTime();
  
        var onerror = function() {
          Ajax.clear(settings);
          var endTime = getTime();
          settings.error(((endTime - startTime) > settings.timeout/2) ? 304 : 403);
        };
  
        var onload = function() {
          Ajax.clear(settings);
          settings.load();
        };
  
        var onsuccess = function() {
          settings.afterSuccess = true;
          settings.success.apply(null, arguments);
        };
  
        script.onerror = onerror;
        script.onload = script.onreadystatechange = function(eventLoad) {
          if (!script.readyState || /loaded|complete/.test(script.readyState)) {
            if (settings.afterSuccess) {
              onload();
            } else {
              onerror();
            }
          }
        };
  
        if (settings.beforeOpen) { settings.beforeOpen({}); }
        if (settings.beforeSend) { settings.beforeSend({}); }
  
        settings.timeoutId = window.setTimeout(onerror, settings.timeout + 2000);
        settings.scriptId = settings.scriptId || getTime();
        settings.afterSuccess = null;
  
        settings.data.callback = settings.scriptId + "_onmessage_" + getTime();
        window[settings.data.callback] = onsuccess;
  
        script.setAttribute("src", addParamsToUrl(settings.url, extend({}, settings.data, currentTimestampParam())));
        script.setAttribute("async", "async");
        script.setAttribute("id", settings.scriptId);
  
        // Use insertBefore instead of appendChild to circumvent an IE6 bug.
        head.insertBefore(script, head.firstChild);
        return settings;
      },
  
      load : function(settings) {
        return Ajax._send(settings, false);
      },
  
      post : function(settings) {
        return Ajax._send(settings, true);
      }
    };
  
    var escapeText = function(text) {
      return (text) ? window.encodeURIComponent(text) : '';
    };
  
    var unescapeText = function(text) {
      return (text) ? window.decodeURIComponent(text) : '';
    };
  
    Utils.parseMessage = function(messageText, keys) {
      var msg = messageText;
      if (isString(messageText)) {
        msg = Utils.parseJSON(messageText);
      }
  
      var message = {
          id     : msg[keys.jsonIdKey],
          channel: msg[keys.jsonChannelKey],
          text   : isString(msg[keys.jsonTextKey]) ? unescapeText(msg[keys.jsonTextKey]) : msg[keys.jsonTextKey],
          tag    : msg[keys.jsonTagKey],
          time   : msg[keys.jsonTimeKey],
          eventid: msg[keys.jsonEventIdKey] || ""
      };
  
      return message;
    };
  
    var getBacktrack = function(options) {
      return (options.backtrack) ? ".b" + Number(options.backtrack) : "";
    };
  
    var getChannelsPath = function(channels, withBacktrack) {
      var path = '';
      for (var channelName in channels) {
        if (!channels.hasOwnProperty || channels.hasOwnProperty(channelName)) {
          path += "/" + channelName + (withBacktrack ? getBacktrack(channels[channelName]) : "");
        }
      }
      return path;
    };
  
    var getSubscriberUrl = function(pushstream, prefix, extraParams, withBacktrack) {
      var websocket = pushstream.wrapper.type === WebSocketWrapper.TYPE;
      var useSSL = pushstream.useSSL;
      var port = normalizePort(useSSL, pushstream.port);
      var url = (websocket) ? ((useSSL) ? "wss://" : "ws://") : ((useSSL) ? "https://" : "http://");
      url += pushstream.host;
      url += (port ? (":" + port) : "");
      url += prefix;
  
      var channels = getChannelsPath(pushstream.channels, withBacktrack);
      if (pushstream.channelsByArgument) {
        var channelParam = {};
        channelParam[pushstream.channelsArgument] = channels.substring(1);
        extraParams = extend({}, extraParams, channelParam);
      } else {
        url += channels;
      }
      return addParamsToUrl(url, extraParams);
    };
  
    var getPublisherUrl = function(pushstream) {
      var port = normalizePort(pushstream.useSSL, pushstream.port);
      var url = (pushstream.useSSL) ? "https://" : "http://";
      url += pushstream.host;
      url += (port ? (":" + port) : "");
      url += pushstream.urlPrefixPublisher;
      url += "?id=" + getChannelsPath(pushstream.channels, false);
  
      return url;
    };
  
    Utils.extract_xss_domain = function(domain) {
      // if domain is an ip address return it, else return ate least the last two parts of it
      if (domain.match(/^(\d{1,3}\.){3}\d{1,3}$/)) {
        return domain;
      }
  
      var domainParts = domain.split('.');
      // if the domain ends with 3 chars or 2 chars preceded by more than 4 chars,
      // we can keep only 2 parts, else we have to keep at least 3 (or all domain name)
      var keepNumber = Math.max(domainParts.length - 1, (domain.match(/(\w{4,}\.\w{2}|\.\w{3,})$/) ? 2 : 3));
  
      return domainParts.slice(-1 * keepNumber).join('.');
    };
  
    var normalizePort = function (useSSL, port) {
      port = Number(port || (useSSL ? 443 : 80));
      return ((!useSSL && port === 80) || (useSSL && port === 443)) ? "" : port;
    };
  
    Utils.isCrossDomainUrl = function(url) {
      if (!url) {
        return false;
      }
  
      var parser = document.createElement('a');
      parser.href = url;
  
      var srcPort = normalizePort(window.location.protocol === "https:", window.location.port);
      var dstPort = normalizePort(parser.protocol === "https:", parser.port);
  
      return (window.location.protocol !== parser.protocol) ||
             (window.location.hostname !== parser.hostname) ||
             (srcPort !== dstPort);
    };
  
    var linker = function(method, instance) {
      return function() {
        return method.apply(instance, arguments);
      };
    };
  
    var clearTimer = function(timer) {
      if (timer) {
        window.clearTimeout(timer);
      }
      return null;
    };
  
    /* common callbacks */
    var onmessageCallback = function(event) {
      Log4js.info("[" + this.type + "] message received", arguments);
      var message = Utils.parseMessage(event.data, this.pushstream);
      if (message.tag) { this.pushstream._etag = message.tag; }
      if (message.time) { this.pushstream._lastModified = message.time; }
      if (message.eventid) { this.pushstream._lastEventId = message.eventid; }
      this.pushstream._onmessage(message.text, message.id, message.channel, message.eventid, true, message.time);
    };
  
    var onopenCallback = function() {
      this.pushstream._onopen();
      Log4js.info("[" + this.type + "] connection opened");
    };
  
    var onerrorCallback = function(event) {
      Log4js.info("[" + this.type + "] error (disconnected by server):", event);
      if ((this.pushstream.readyState === PushStream.OPEN) &&
          (this.type === EventSourceWrapper.TYPE) &&
          (event.type === 'error') &&
          (this.connection.readyState === window.EventSource.CONNECTING)) {
        // EventSource already has a reconnection function using the last-event-id header
        return;
      }
      this._closeCurrentConnection();
      this.pushstream._onerror({type: ((event && ((event.type === "load") || ((event.type === "close") && (event.code === 1006)))) || (this.pushstream.readyState === PushStream.CONNECTING)) ? "load" : "timeout"});
    };
  
    /* wrappers */
  
    var WebSocketWrapper = function(pushstream) {
      if (!window.WebSocket && !window.MozWebSocket) { throw "WebSocket not supported"; }
      this.type = WebSocketWrapper.TYPE;
      this.pushstream = pushstream;
      this.connection = null;
    };
  
    WebSocketWrapper.TYPE = "WebSocket";
  
    WebSocketWrapper.prototype = {
      connect: function() {
        this._closeCurrentConnection();
        var params = extend({}, this.pushstream.extraParams(), currentTimestampParam(), getControlParams(this.pushstream));
        var url = getSubscriberUrl(this.pushstream, this.pushstream.urlPrefixWebsocket, params, !this.pushstream._useControlArguments());
        this.connection = (window.WebSocket) ? new window.WebSocket(url) : new window.MozWebSocket(url);
        this.connection.onerror   = linker(onerrorCallback, this);
        this.connection.onclose   = linker(onerrorCallback, this);
        this.connection.onopen    = linker(onopenCallback, this);
        this.connection.onmessage = linker(onmessageCallback, this);
        Log4js.info("[WebSocket] connecting to:", url);
      },
  
      disconnect: function() {
        if (this.connection) {
          Log4js.debug("[WebSocket] closing connection to:", this.connection.url);
          this.connection.onclose = null;
          this._closeCurrentConnection();
          this.pushstream._onclose();
        }
      },
  
      _closeCurrentConnection: function() {
        if (this.connection) {
          try { this.connection.close(); } catch (e) { /* ignore error on closing */ }
          this.connection = null;
        }
      },
  
      sendMessage: function(message) {
        if (this.connection) { this.connection.send(message); }
      }
    };
  
    var EventSourceWrapper = function(pushstream) {
      if (!window.EventSource) { throw "EventSource not supported"; }
      this.type = EventSourceWrapper.TYPE;
      this.pushstream = pushstream;
      this.connection = null;
    };
  
    EventSourceWrapper.TYPE = "EventSource";
  
    EventSourceWrapper.prototype = {
      connect: function() {
        this._closeCurrentConnection();
        var params = extend({}, this.pushstream.extraParams(), currentTimestampParam(), getControlParams(this.pushstream));
        var url = getSubscriberUrl(this.pushstream, this.pushstream.urlPrefixEventsource, params, !this.pushstream._useControlArguments());
        this.connection = new window.EventSource(url);
        this.connection.onerror   = linker(onerrorCallback, this);
        this.connection.onopen    = linker(onopenCallback, this);
        this.connection.onmessage = linker(onmessageCallback, this);
        Log4js.info("[EventSource] connecting to:", url);
      },
  
      disconnect: function() {
        if (this.connection) {
          Log4js.debug("[EventSource] closing connection to:", this.connection.url);
          this.connection.onclose = null;
          this._closeCurrentConnection();
          this.pushstream._onclose();
        }
      },
  
      _closeCurrentConnection: function() {
        if (this.connection) {
          try { this.connection.close(); } catch (e) { /* ignore error on closing */ }
          this.connection = null;
        }
      }
    };
  
    var StreamWrapper = function(pushstream) {
      this.type = StreamWrapper.TYPE;
      this.pushstream = pushstream;
      this.connection = null;
      this.url = null;
      this.frameloadtimer = null;
      this.pingtimer = null;
      this.iframeId = "PushStreamManager_" + pushstream.id;
    };
  
    StreamWrapper.TYPE = "Stream";
  
    StreamWrapper.prototype = {
      connect: function() {
        this._closeCurrentConnection();
        var domain = Utils.extract_xss_domain(this.pushstream.host);
        try {
          document.domain = domain;
        } catch(e) {
          Log4js.error("[Stream] (warning) problem setting document.domain = " + domain + " (OBS: IE8 does not support set IP numbers as domain)");
        }
        var params = extend({}, this.pushstream.extraParams(), currentTimestampParam(), {"streamid": this.pushstream.id}, getControlParams(this.pushstream));
        this.url = getSubscriberUrl(this.pushstream, this.pushstream.urlPrefixStream, params, !this.pushstream._useControlArguments());
        Log4js.debug("[Stream] connecting to:", this.url);
        this.loadFrame(this.url);
      },
  
      disconnect: function() {
        if (this.connection) {
          Log4js.debug("[Stream] closing connection to:", this.url);
          this._closeCurrentConnection();
          this.pushstream._onclose();
        }
      },
  
      _clear_iframe: function() {
        var oldIframe = document.getElementById(this.iframeId);
        if (oldIframe) {
          oldIframe.onload = null;
          oldIframe.src = "about:blank";
          if (oldIframe.parentNode) { oldIframe.parentNode.removeChild(oldIframe); }
        }
      },
  
      _closeCurrentConnection: function() {
        this._clear_iframe();
        if (this.connection) {
          this.pingtimer = clearTimer(this.pingtimer);
          this.frameloadtimer = clearTimer(this.frameloadtimer);
          this.connection = null;
          this.transferDoc = null;
          if (typeof window.CollectGarbage === 'function') { window.CollectGarbage(); }
        }
      },
  
      loadFrame: function(url) {
        this._clear_iframe();
  
        var ifr = null;
        if ("ActiveXObject" in window) {
          var transferDoc = new window.ActiveXObject("htmlfile");
          transferDoc.open();
          transferDoc.write("\x3C" + "html" + "\x3E\x3C" + "script" + "\x3E" + "document.domain='" + document.domain + "';\x3C" + "/script" + "\x3E");
          transferDoc.write("\x3C" + "body" + "\x3E\x3C" + "iframe id='" + this.iframeId + "' src='" + url + "'\x3E\x3C" + "/iframe" + "\x3E\x3C" + "/body" + "\x3E\x3C" + "/html" + "\x3E");
          transferDoc.parentWindow.PushStream = PushStream;
          transferDoc.close();
          ifr = transferDoc.getElementById(this.iframeId);
          this.transferDoc = transferDoc;
        } else {
          ifr = document.createElement("IFRAME");
          ifr.style.width = "1px";
          ifr.style.height = "1px";
          ifr.style.border = "none";
          ifr.style.position = "absolute";
          ifr.style.top = "-10px";
          ifr.style.marginTop = "-10px";
          ifr.style.zIndex = "-20";
          ifr.PushStream = PushStream;
          document.body.appendChild(ifr);
          ifr.setAttribute("src", url);
          ifr.setAttribute("id", this.iframeId);
        }
  
        ifr.onload = linker(onerrorCallback, this);
        this.connection = ifr;
        this.frameloadtimer = window.setTimeout(linker(onerrorCallback, this), this.pushstream.timeout);
      },
  
      register: function(iframeWindow) {
        this.frameloadtimer = clearTimer(this.frameloadtimer);
        iframeWindow.p = linker(this.process, this);
        this.connection.onload = linker(this._onframeloaded, this);
        this.pushstream._onopen();
        this.setPingTimer();
        Log4js.info("[Stream] frame registered");
      },
  
      process: function(id, channel, text, eventid, time, tag) {
        this.pingtimer = clearTimer(this.pingtimer);
        Log4js.info("[Stream] message received", arguments);
        if (id !== -1) {
          if (tag) { this.pushstream._etag = tag; }
          if (time) { this.pushstream._lastModified = time; }
          if (eventid) { this.pushstream._lastEventId = eventid; }
        }
        this.pushstream._onmessage(unescapeText(text), id, channel, eventid || "", true, time);
        this.setPingTimer();
      },
  
      _onframeloaded: function() {
        Log4js.info("[Stream] frame loaded (disconnected by server)");
        this.pushstream._onerror({type: "timeout"});
        this.connection.onload = null;
        this.disconnect();
      },
  
      setPingTimer: function() {
        if (this.pingtimer) { clearTimer(this.pingtimer); }
        this.pingtimer = window.setTimeout(linker(onerrorCallback, this), this.pushstream.pingtimeout);
      }
    };
  
    var LongPollingWrapper = function(pushstream) {
      this.type = LongPollingWrapper.TYPE;
      this.pushstream = pushstream;
      this.connection = null;
      this.opentimer = null;
      this.messagesQueue = [];
      this._linkedInternalListen = linker(this._internalListen, this);
      this.xhrSettings = {
          timeout: this.pushstream.timeout,
          data: {},
          url: null,
          success: linker(this.onmessage, this),
          error: linker(this.onerror, this),
          load: linker(this.onload, this),
          beforeSend: linker(this.beforeSend, this),
          afterReceive: linker(this.afterReceive, this)
      };
    };
  
    LongPollingWrapper.TYPE = "LongPolling";
  
    LongPollingWrapper.prototype = {
      connect: function() {
        this.messagesQueue = [];
        this._closeCurrentConnection();
        this.urlWithBacktrack = getSubscriberUrl(this.pushstream, this.pushstream.urlPrefixLongpolling, {}, true);
        this.urlWithoutBacktrack = getSubscriberUrl(this.pushstream, this.pushstream.urlPrefixLongpolling, {}, false);
        this.xhrSettings.url = this.urlWithBacktrack;
        this.useJSONP = this.pushstream._crossDomain || this.pushstream.useJSONP;
        this.xhrSettings.scriptId = "PushStreamManager_" + this.pushstream.id;
        if (this.useJSONP) {
          this.pushstream.messagesControlByArgument = true;
        }
        this._listen();
        this.opentimer = window.setTimeout(linker(onopenCallback, this), 150);
        Log4js.info("[LongPolling] connecting to:", this.xhrSettings.url);
      },
  
      _listen: function() {
        if (this._internalListenTimeout) { clearTimer(this._internalListenTimeout); }
        this._internalListenTimeout = window.setTimeout(this._linkedInternalListen, 100);
      },
  
      _internalListen: function() {
        if (this.pushstream._keepConnected) {
          this.xhrSettings.url = this.pushstream._useControlArguments() ? this.urlWithoutBacktrack : this.urlWithBacktrack;
          this.xhrSettings.data = extend({}, this.pushstream.extraParams(), this.xhrSettings.data, getControlParams(this.pushstream));
          if (this.useJSONP) {
            this.connection = Ajax.jsonp(this.xhrSettings);
          } else if (!this.connection) {
            this.connection = Ajax.load(this.xhrSettings);
          }
        }
      },
  
      disconnect: function() {
        if (this.connection) {
          Log4js.debug("[LongPolling] closing connection to:", this.xhrSettings.url);
          this._closeCurrentConnection();
          this.pushstream._onclose();
        }
      },
  
      _closeCurrentConnection: function() {
        this.opentimer = clearTimer(this.opentimer);
        if (this.connection) {
          try { this.connection.abort(); } catch (e) {
            try { Ajax.clear(this.connection); } catch (e1) { /* ignore error on closing */ }
          }
          this.connection = null;
          this.xhrSettings.url = null;
        }
      },
  
      beforeSend: function(xhr) {
        if (!this.pushstream.messagesControlByArgument) {
          xhr.setRequestHeader("If-None-Match", this.pushstream._etag);
          xhr.setRequestHeader("If-Modified-Since", this.pushstream._lastModified);
        }
      },
  
      afterReceive: function(xhr) {
        if (!this.pushstream.messagesControlByArgument) {
          this.pushstream._etag = xhr.getResponseHeader('Etag');
          this.pushstream._lastModified = xhr.getResponseHeader('Last-Modified');
        }
        this.connection = null;
      },
  
      onerror: function(status) {
        this._closeCurrentConnection();
        if (this.pushstream._keepConnected) { /* abort(), called by disconnect(), call this callback, but should be ignored */
          if (status === 304) {
            this._listen();
          } else {
            Log4js.info("[LongPolling] error (disconnected by server):", status);
            this.pushstream._onerror({type: ((status === 403) || (this.pushstream.readyState === PushStream.CONNECTING)) ? "load" : "timeout"});
          }
        }
      },
  
      onload: function() {
        this._listen();
      },
  
      onmessage: function(responseText) {
        if (this._internalListenTimeout) { clearTimer(this._internalListenTimeout); }
        Log4js.info("[LongPolling] message received", responseText);
        var lastMessage = null;
        var messages = isArray(responseText) ? responseText : responseText.replace(/\}\{/g, "}\r\n{").split("\r\n");
        for (var i = 0; i < messages.length; i++) {
          if (messages[i]) {
            lastMessage = Utils.parseMessage(messages[i], this.pushstream);
            this.messagesQueue.push(lastMessage);
            if (this.pushstream.messagesControlByArgument && lastMessage.time) {
              this.pushstream._etag = lastMessage.tag;
              this.pushstream._lastModified = lastMessage.time;
            }
          }
        }
  
        this._listen();
  
        while (this.messagesQueue.length > 0) {
          var message = this.messagesQueue.shift();
          this.pushstream._onmessage(message.text, message.id, message.channel, message.eventid, (this.messagesQueue.length === 0), message.time);
        }
      }
    };
  
    /* mains class */
  
    var PushStreamManager = [];
  
    var PushStream = function(settings) {
      settings = settings || {};
  
      this.id = PushStreamManager.push(this) - 1;
  
      this.useSSL = settings.useSSL || false;
      this.host = settings.host || window.location.hostname;
      this.port = Number(settings.port || (this.useSSL ? 443 : 80));
  
      this.timeout = settings.timeout || 30000;
      this.pingtimeout = settings.pingtimeout || 30000;
      this.reconnectOnTimeoutInterval = settings.reconnectOnTimeoutInterval || 3000;
      this.reconnectOnChannelUnavailableInterval = settings.reconnectOnChannelUnavailableInterval || 60000;
      this.autoReconnect = (settings.autoReconnect !== false);
  
      this.lastEventId = settings.lastEventId || null;
      this.messagesPublishedAfter = settings.messagesPublishedAfter;
      this.messagesControlByArgument = settings.messagesControlByArgument || false;
      this.tagArgument   = settings.tagArgument  || 'tag';
      this.timeArgument  = settings.timeArgument || 'time';
      this.eventIdArgument  = settings.eventIdArgument || 'eventid';
      this.useJSONP      = settings.useJSONP     || false;
  
      this._reconnecttimer = null;
      this._etag = 0;
      this._lastModified = null;
      this._lastEventId = null;
  
      this.urlPrefixPublisher   = settings.urlPrefixPublisher   || '/pub';
      this.urlPrefixStream      = settings.urlPrefixStream      || '/sub';
      this.urlPrefixEventsource = settings.urlPrefixEventsource || '/ev';
      this.urlPrefixLongpolling = settings.urlPrefixLongpolling || '/lp';
      this.urlPrefixWebsocket   = settings.urlPrefixWebsocket   || '/ws';
  
      this.jsonIdKey      = settings.jsonIdKey      || 'id';
      this.jsonChannelKey = settings.jsonChannelKey || 'channel';
      this.jsonTextKey    = settings.jsonTextKey    || 'text';
      this.jsonTagKey     = settings.jsonTagKey     || 'tag';
      this.jsonTimeKey    = settings.jsonTimeKey    || 'time';
      this.jsonEventIdKey = settings.jsonEventIdKey || 'eventid';
  
      this.modes = (settings.modes || 'eventsource|websocket|stream|longpolling').split('|');
      this.wrappers = [];
      this.wrapper = null;
  
      this.onchanneldeleted = settings.onchanneldeleted || null;
      this.onmessage = settings.onmessage || null;
      this.onerror = settings.onerror || null;
      this.onstatuschange = settings.onstatuschange || null;
      this.extraParams    = settings.extraParams    || function() { return {}; };
  
      this.channels = {};
      this.channelsCount = 0;
      this.channelsByArgument   = settings.channelsByArgument   || false;
      this.channelsArgument     = settings.channelsArgument     || 'channels';
  
      this._crossDomain = Utils.isCrossDomainUrl(getPublisherUrl(this));
  
      for (var i = 0; i < this.modes.length; i++) {
        try {
          var wrapper = null;
          switch (this.modes[i]) {
          case "websocket"  : wrapper = new WebSocketWrapper(this);   break;
          case "eventsource": wrapper = new EventSourceWrapper(this); break;
          case "longpolling": wrapper = new LongPollingWrapper(this); break;
          case "stream"     : wrapper = new StreamWrapper(this);      break;
          }
          this.wrappers[this.wrappers.length] = wrapper;
        } catch(e) { Log4js.info(e); }
      }
  
      this.readyState = 0;
    };
  
    /* constants */
    PushStream.LOG_LEVEL = 'error'; /* debug, info, error */
    PushStream.LOG_OUTPUT_ELEMENT_ID = 'Log4jsLogOutput';
  
    /* status codes */
    PushStream.CLOSED        = 0;
    PushStream.CONNECTING    = 1;
    PushStream.OPEN          = 2;
  
    /* main code */
    PushStream.prototype = {
      addChannel: function(channel, options) {
        if (escapeText(channel) !== channel) {
          throw "Invalid channel name! Channel has to be a set of [a-zA-Z0-9]";
        }
        Log4js.debug("entering addChannel");
        if (typeof(this.channels[channel]) !== "undefined") { throw "Cannot add channel " + channel + ": already subscribed"; }
        options = options || {};
        Log4js.info("adding channel", channel, options);
        this.channels[channel] = options;
        this.channelsCount++;
        if (this.readyState !== PushStream.CLOSED) { this.connect(); }
        Log4js.debug("leaving addChannel");
      },
  
      removeChannel: function(channel) {
        if (this.channels[channel]) {
          Log4js.info("removing channel", channel);
          delete this.channels[channel];
          this.channelsCount--;
        }
      },
  
      removeAllChannels: function() {
        Log4js.info("removing all channels");
        this.channels = {};
        this.channelsCount = 0;
      },
  
      _setState: function(state) {
        if (this.readyState !== state) {
          Log4js.info("status changed", state);
          this.readyState = state;
          if (this.onstatuschange) {
            this.onstatuschange(this.readyState);
          }
        }
      },
  
      connect: function() {
        Log4js.debug("entering connect");
        if (!this.host)                 { throw "PushStream host not specified"; }
        if (isNaN(this.port))           { throw "PushStream port not specified"; }
        if (!this.channelsCount)        { throw "No channels specified"; }
        if (this.wrappers.length === 0) { throw "No available support for this browser"; }
  
        this._keepConnected = true;
        this._lastUsedMode = 0;
        this._connect();
  
        Log4js.debug("leaving connect");
      },
  
      disconnect: function() {
        Log4js.debug("entering disconnect");
        this._keepConnected = false;
        this._disconnect();
        this._setState(PushStream.CLOSED);
        Log4js.info("disconnected");
        Log4js.debug("leaving disconnect");
      },
  
      _useControlArguments :function() {
        return this.messagesControlByArgument && ((this._lastModified !== null) || (this._lastEventId !== null));
      },
  
      _connect: function() {
        if (this._lastEventId === null) {
          this._lastEventId = this.lastEventId;
        }
        if (this._lastModified === null) {
          var date = this.messagesPublishedAfter;
          if (!isDate(date)) {
            var messagesPublishedAfter = Number(this.messagesPublishedAfter);
            if (messagesPublishedAfter > 0) {
              date = new Date();
              date.setTime(date.getTime() - (messagesPublishedAfter * 1000));
            } else if (messagesPublishedAfter < 0) {
              date = new Date(0);
            }
          }
  
          if (isDate(date)) {
            this._lastModified = Utils.dateToUTCString(date);
          }
        }
  
        this._disconnect();
        this._setState(PushStream.CONNECTING);
        this.wrapper = this.wrappers[this._lastUsedMode++ % this.wrappers.length];
  
        try {
          this.wrapper.connect();
        } catch (e) {
          //each wrapper has a cleanup routine at disconnect method
          if (this.wrapper) {
            this.wrapper.disconnect();
          }
        }
      },
  
      _disconnect: function() {
        this._reconnecttimer = clearTimer(this._reconnecttimer);
        if (this.wrapper) {
          this.wrapper.disconnect();
        }
      },
  
      _onopen: function() {
        this._reconnecttimer = clearTimer(this._reconnecttimer);
        this._setState(PushStream.OPEN);
        if (this._lastUsedMode > 0) {
          this._lastUsedMode--; //use same mode on next connection
        }
      },
  
      _onclose: function() {
        this._reconnecttimer = clearTimer(this._reconnecttimer);
        this._setState(PushStream.CLOSED);
        this._reconnect(this.reconnectOnTimeoutInterval);
      },
  
      _onmessage: function(text, id, channel, eventid, isLastMessageFromBatch, time) {
        Log4js.debug("message", text, id, channel, eventid, isLastMessageFromBatch, time);
        if (id === -2) {
          if (this.onchanneldeleted) { this.onchanneldeleted(channel); }
        } else if (id > 0) {
          if (this.onmessage) { this.onmessage(text, id, channel, eventid, isLastMessageFromBatch, time); }
        }
      },
  
      _onerror: function(error) {
        this._setState(PushStream.CLOSED);
        this._reconnect((error.type === "timeout") ? this.reconnectOnTimeoutInterval : this.reconnectOnChannelUnavailableInterval);
        if (this.onerror) { this.onerror(error); }
      },
  
      _reconnect: function(timeout) {
        if (this.autoReconnect && this._keepConnected && !this._reconnecttimer && (this.readyState !== PushStream.CONNECTING)) {
          Log4js.info("trying to reconnect in", timeout);
          this._reconnecttimer = window.setTimeout(linker(this._connect, this), timeout);
        }
      },
  
      sendMessage: function(message, successCallback, errorCallback) {
        message = escapeText(message);
        if (this.wrapper.type === WebSocketWrapper.TYPE) {
          this.wrapper.sendMessage(message);
          if (successCallback) { successCallback(); }
        } else {
          Ajax.post({url: getPublisherUrl(this), data: message, success: successCallback, error: errorCallback, crossDomain: this._crossDomain});
        }
      }
    };
  
    PushStream.sendMessage = function(url, message, successCallback, errorCallback) {
      Ajax.post({url: url, data: escapeText(message), success: successCallback, error: errorCallback});
    };
  
    // to make server header template more clear, it calls register and
    // by a url parameter we find the stream wrapper instance
    PushStream.register = function(iframe) {
      var matcher = iframe.window.location.href.match(/streamid=([0-9]*)&?/);
      if (matcher[1] && PushStreamManager[matcher[1]]) {
        PushStreamManager[matcher[1]].wrapper.register(iframe);
      }
    };
  
    PushStream.unload = function() {
      for (var i = 0; i < PushStreamManager.length; i++) {
        try { PushStreamManager[i].disconnect(); } catch(e){}
      }
    };
  
    /* make class public */
    window.PushStream = PushStream;
    window.PushStreamManager = PushStreamManager;
  
    if (window.attachEvent) { window.attachEvent("onunload", PushStream.unload); }
    if (window.addEventListener) { window.addEventListener.call(window, "unload", PushStream.unload, false); }
  
})(window, document);
  
// shhgit
document.addEventListener('DOMContentLoaded', function(event) {
    window.connection = null;
    window.timeout = null;

    var settings = {
        activeSignatures: [],
        burger: document.getElementById('burger'),
        connectionStats: document.getElementById('connection-status'),
        interestingFiles: document.getElementById('setting-interesting-files'),
        highEntropyStrings: document.getElementById('setting-high-entropy-strings'),
        notifications: document.getElementById('setting-notifications'),
        matchesCount: document.getElementById('matches-count').getElementsByTagName('span')[0],
        filtersClear: document.getElementById('filters-clear'),
        filtersCount: document.getElementById('filters-count').getElementsByTagName('span')[0]
    };
    const slugify = (value) => value.toLowerCase().replace(/[^a-z0-9 -]/g, '').replace(/\s+/g, '-').replace(/-+/g, '-');
    const getFileUrl = (data) => {
        if (!data.Url.substr(-4) === '.git') return data.Url;

        var source = getSource(data.Source);
        var prefix = source.icon == 'bitbucket' ? 'src' : 'blob';
        
        return `${data.Url.substr(0, data.Url.indexOf('.git'))}/${prefix}/master${data.File}`;
    };
    const getSource = (source) => {
        switch (source) {
            case 0: return {icon: 'github', name: 'GitHub'};
            case 1: return {icon: 'github-square', name: 'Gist'};
            case 2: return {icon: 'bitbucket', name: 'BitBucket'};
            case 3: return {icon: 'gitlab', name: 'GitLab'};
        }
    };
    const getIssueUrl = (data) => {
        var root = data.Url.substr(0, data.Url.indexOf('.git'));
        var title = encodeURIComponent(`Exposed ${data.Signature}`);
        var description = encodeURIComponent(`Potential security breach. See ${data.File}`);
        var source = getSource(data.Source);

        switch (source.icon) {
            case 'github': return `${root}/issues/new?title=${title}&body=${description}`;
            case 'gitlab': return `${root}/issues/new?issue[title]=${title}&issue[description]=${description}`;
        }
    };
    const sort = (list) => {      
        signatures = list.getElementsByTagName("li");
        Array.from(signatures)
            .sort((a, b) => parseInt(b.getElementsByClassName('menu-item')[0].getAttribute('data-badge') || 0) - parseInt(a.getElementsByClassName('menu-item')[0].getAttribute('data-badge') || 0))
            .forEach(li => list.appendChild(li));
    };
    const updateStatus = (text, cls) => {
        settings.connectionStats.classList.remove('is-info', 'is-success', 'is-warning', 'is-danger');
        settings.connectionStats.classList.add(cls);
        settings.connectionStats.textContent = text;
    };
    const filterSignature = (signature) => {
        var state = settings.activeSignatures.includes(signature.id);
        signature.classList.toggle('is-active');

        if (!state) settings.activeSignatures.push(signature.id);

        Array.from(document.getElementsByClassName('log')).forEach(log => log.style.display = 'none');
        settings.activeSignatures.forEach((signatureId) => {
            Array.from(document.getElementsByClassName(signatureId)).forEach(log => {
                if (state && signatureId == signature.id) return;
                if (!settings.interestingFiles.checked && log.classList.contains('is-interesting-file')) return;
                if (!settings.highEntropyStrings.checked && log.classList.contains('is-high-entropy-string')) return;

                log.style.display = '';
            });
        });
        
        if (state) {
            settings.activeSignatures.splice(settings.activeSignatures.indexOf(signature.id), 1);
            var anyActive = (settings.activeSignatures.length > 0);
            Array.from(document.getElementsByClassName(anyActive ? signature.id : 'log')).forEach(log => {
                if (!settings.interestingFiles.checked && log.classList.contains('is-interesting-file')) return;
                if (!settings.highEntropyStrings.checked && log.classList.contains('is-high-entropy-string')) return;

                log.style.display = anyActive ? 'none' : '';
            });
        }

        settings.filtersCount.textContent = `${settings.activeSignatures.length} filters`;
    };
    const processEvent = (data) => {
        var eventId = CryptoJS.MD5(data.File + '-' + data.Signature + '-' + (data.Matches ? data.Matches.join('') : '0')).toString();
        if (document.getElementById(eventId)) return; // duplicate

        var sigId = slugify(data.Signature);
        var matchesCount = data.Matches ? data.Matches.length : 1;
        var sigMenuItem = document.getElementById(sigId).getElementsByClassName('menu-item')[0];
        var source = getSource(data.Source)
        sigMenuItem.setAttribute('data-badge', parseInt(sigMenuItem.getAttribute('data-badge') || 0) + matchesCount);
        sort(document.getElementById('signatures'));

        var row = document.getElementById('messages').insertRow(0);
        row.classList.add('log', sigId);
        row.id = eventId;
        row.insertCell(0).innerHTML = `<td class="source"><span class="icon" title="${source.name}"><i class="fab fa-lg fa-${source.icon}"></i></span></td>`;
        row.insertCell(1).innerHTML = `<td class="found"><span class="datetime" title="${new Date().toLocaleString}">${new Date().toLocaleTimeString()}</span></td>`;
        row.insertCell(2).innerHTML = `<td class="signature-name"><strong>${data.Signature}</strong>${source.icon != 'bitbucket' ? `<a href="${getIssueUrl(data)}" title="Raise an issue" target="_blank" onclick="event.stopPropagation();"><span class="icon is-dark"><i class="fas fa-flag"></i></span></a>` : ''}</td>`;
        row.insertCell(3).innerHTML = `<td class="matches"><div>${data.Matches ? "<pre>" + data.Matches.join('<br />') + "</pre>" : '<em>&mdash;</em>'}</div></td>`;
        row.insertCell(4).innerHTML = `<td class="file-url"><a href="${getFileUrl(data)}" target="_blank">${data.File}</a></td>`;
        row.insertCell(5).innerHTML = `<td class="stars">${data.Stars}</td>`
        row.addEventListener('click', (event) => {
            event.preventDefault();
            window.open(getFileUrl(data), '_blank');
        });
        settings.matchesCount.textContent = `${document.getElementsByClassName('log').length} matches`;

        if (!data.Matches) {
            row.classList.add('is-interesting-file')
            if (!settings.interestingFiles.checked) row.style.display = 'none';
        }

        if (data.Signature === "High entropy string") {
            row.classList.add('is-high-entropy-string')
            if (!settings.highEntropyStrings.checked) row.style.display = 'none';
        }

        if (settings.activeSignatures.length > 0 && !settings.activeSignatures.includes(sigId)) row.style.display = 'none';
        if (settings.notifications.checked) notifyFinding(data.Signature, data.Matches ? data.Matches.join(', ') : data.File);
    };
    const listenForEvents = () => {
        window.connection = new PushStream({
            host: 'localhost',
            port: 8080,
            urlPrefixEventsource: '/events',
            useSSL: false,
            modes: 'eventsource',
            messagesPublishedAfter: 100,
            messagesControlByArgument: true
        });

        window.connection.onerror = (e) => {
          if (confirm("Error connecting to shhgit. Reload to retry?")) window.location.reload();
        };
        window.connection.onstatuschange = (e) => {
            if (e == PushStream.OPEN) updateStatus('Connected', 'is-success');
            else if (e == PushStream.CONNECTING) updateStatus('Syncing...', 'is-info');
        };
        window.connection.onmessage = (text, id, channel, eventid, isLast, time) => {
            if (document.getElementById('loading')) document.getElementById('loading').remove();

            processEvent(text);
        };

        window.connection.addChannel('shhgit');
        window.connection.connect();
    };
    const notifyFinding = (title, message) => {
        if (Notification.permission === "granted") {
            new Notification(title, {
                'icon': '/logo.png',
                'body': message
            });
        }
    };

    (() => {
        burger.addEventListener('click', () => {
            const target = burger.dataset.target;
            const $target = document.getElementById(target);
    
            burger.classList.toggle('is-active');
            $target.classList.toggle('is-active');
        });

        settings.interestingFiles.addEventListener('change', (event) => {
            Array.from(document.getElementsByClassName('is-interesting-file'))
                .forEach(log => {
                    log.style.display = event.target.checked ? '' : 'none'
                });
        });

        settings.highEntropyStrings.addEventListener('change', (event) => {
            Array.from(document.getElementsByClassName('is-high-entropy-string'))
                .forEach(log => {
                    log.style.display = event.target.checked ? '' : 'none'
                });
        });
    
        settings.notifications.addEventListener('change', (event) => {
            Notification.requestPermission().then((permission) => {
                if (permission !== "granted") {
                    settings.notifications.checked = false;
                    settings.notifications.disabled = true;
                }
            });
        });

        settings.filtersClear.addEventListener('click', (event) => {
            settings.activeSignatures = [];
            settings.filtersCount.textContent = "0 filters";

            Array.from(document.querySelectorAll('#signatures li.is-active')).forEach(log => log.classList.remove('is-active'));
            Array.from(document.getElementsByClassName('log')).forEach(log => log.style.display = '');
        });

        fetch(`/signatures.json`)
            .then((resp) => resp.json())
            .then(signatures => {
                signatures.forEach(signature => {
                    var li = document.createElement('li');
                    li.id = slugify(signature)
                    li.innerHTML = `<a href="#" class="menu-item" title="${signature}">${signature}</a>`;
                    li.addEventListener('click', (event) => {
                        event.preventDefault();
                        filterSignature(li);
                    });

                    document.getElementById('signatures').appendChild(li);
                });

                listenForEvents();
            })
            .catch((err) => {
                alert('Failed to retrieve signatures! Reloading...')
                window.location.reload();
            });
    })();
});