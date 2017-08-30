import 'dart:html';

import 'package:angular2/angular2.dart';
import 'package:gitpods/user.dart';
import 'package:gitpods/user_service.dart';
import 'package:angular2/router.dart';

@Component(
  selector: 'gitpods-user-settings',
  templateUrl: 'user_settings_component.html',
  directives: const [COMMON_DIRECTIVES],
)
class UserSettingsComponent implements OnInit {
  final Router _router;
  final UserService _userService;

  UserSettingsComponent(this._router, this._userService);

  User user;

  @override
  ngOnInit() {
    this._userService.me()
        .then((user) => this.user = user);
  }

  submit(Event event) {
    event.preventDefault();


    this._userService.update(user)
        .then((user) => this._router.navigate(['UserProfile', {'username': user.username}]));
  }
}
