import 'dart:html';

import 'package:angular/angular.dart';
import 'package:angular_forms/angular_forms.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/user.dart';
import 'package:gitpods/user_service.dart';

@Component(
  selector: 'profile',
  templateUrl: 'profile_component.html',
  directives: const [COMMON_DIRECTIVES, formDirectives],
)
class ProfileComponent implements OnInit {
  final Router _router;
  final UserService _userService;

  ProfileComponent(this._router, this._userService);

  User user;

  @override
  void ngOnInit() {
    this._userService.me()
        .then((user) => this.user = user);
  }

  void submit(Event event) {
    event.preventDefault();

    this._userService.update(user)
        .then((user) => this._router.navigate(['/UserProfile', {'username': user.username}]));
  }
}
