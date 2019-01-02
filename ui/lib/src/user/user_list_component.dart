import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/routes.dart';
import 'package:gitpods/src/gravatar_component.dart';
import 'package:gitpods/src/loading_component.dart';
import 'package:gitpods/src/user/user.dart';
import 'package:gitpods/src/user/user_service.dart';

@Component(
  selector: 'gitpods-user-list',
  templateUrl: 'user_list_component.html',
  directives: const [
    coreDirectives,
    routerDirectives,
    LoadingComponent,
    GravatarComponent,
  ],
  providers: const [UserService],
)
class UserListComponent implements OnInit {
  UserListComponent(this._userService);

  final UserService _userService;

  List<User> users;

  @override
  void ngOnInit() {
    _userService
        .list()
        .then((List<User> users) => this.users = users)
        .catchError((e) => print(e.toString()));
  }

  String userProfileUrl(String username) {
    return RoutePaths.userProfile.toUrl(parameters: {
      'username': username,
    });
  }
}
