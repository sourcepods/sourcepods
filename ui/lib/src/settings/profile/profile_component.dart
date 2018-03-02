import 'dart:html';

import 'package:angular/angular.dart';
import 'package:angular_forms/angular_forms.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/src/loading_component.dart';
import 'package:gitpods/src/user/user.dart';
import 'package:gitpods/src/user/user_service.dart';

@Component(
  selector: 'profile',
  templateUrl: 'profile_component.html',
  directives: const [
    COMMON_DIRECTIVES,
    formDirectives,
    LoadingComponent,
  ],
)
class ProfileComponent implements OnInit {
  final Router _router;
  final UserService _userService;

  ProfileComponent(this._router, this._userService);

  bool loading;
  User user;

  @override
  void ngOnInit() {
    this._userService.me().then((user) => this.user = user);
  }

  void submit(Event e) {
    e.preventDefault();
    loading = true;

    this._userService.update(user).then((user) {
      loading = false;
      this._router.navigate(['/UserProfile',{'username': user.username}]);
    });
  }
}
