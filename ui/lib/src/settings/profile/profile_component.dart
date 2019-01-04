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
    coreDirectives,
    formDirectives,
    LoadingComponent,
  ],
  providers: [
    ClassProvider(UserService),
  ]
)
class ProfileComponent implements OnActivate {
  ProfileComponent(this._userService);

  final UserService _userService;

  bool loading;
  User user;

  @override
  void onActivate(RouterState previous, RouterState current) {
    _userService.me().then((User user) => this.user = user);
  }

  void submit(Event e) {
    e.preventDefault();
    loading = true;

    this._userService.update(user).then((user) {
      loading = false;
//      this._router.navigate(['/UserProfile',{'username': user.username}]);
    });
  }
}
