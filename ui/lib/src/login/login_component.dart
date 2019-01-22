import 'dart:html';

import 'package:angular/angular.dart';
import 'package:angular_forms/angular_forms.dart';
import 'package:sourcepods/src/login/login_service.dart';

@Component(
  selector: 'gitpods-selector',
  templateUrl: 'login_component.html',
  styleUrls: const ['login_component.css'],
  directives: const [coreDirectives, formDirectives],
  providers: const [LoginService],
)
class LoginComponent {
  LoginComponent(this._loginService);

  final LoginService _loginService;

  String email;
  String password;

  String error = '';

  void submit(Event e) {
    e.preventDefault();
    this.error = '';

    _loginService
        .login(this.email, this.password)
        .then((String error) => this.error = error);
  }
}
