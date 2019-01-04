import 'dart:async';
import 'dart:convert';
import 'dart:html';

import 'package:angular/angular.dart';
import 'package:http/http.dart';

@Injectable()
class LoginService {
  LoginService(this._http);

  Client _http;

  Future<String> login(String email, String password) async {
    try {
      Map body = {'email': email, 'password': password};
      final resp = await _http.post('/api/authorize', body: json.encode(body));

      if (resp.statusCode == 200) {
        // Reload the page to have the new cookie set.
        window.location.assign('/');
        return '';
      }

      var data = json.decode(resp.body);

      return data['errors'][0]['detail'];
    } catch (e) {
      return e.toString();
    }
  }

  void logout() {
    window.location.assign('/api/sessions/logout');
  }
}
