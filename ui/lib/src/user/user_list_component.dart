import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/src/gravatar_component.dart';
import 'package:gitpods/src/user/user.dart';
import 'package:gitpods/src/user/user_service.dart';

@Component(
  selector: 'gitpods-user-list',
  templateUrl: 'user_list_component.html',
  directives: const [COMMON_DIRECTIVES, ROUTER_DIRECTIVES, GravatarComponent],
  providers: const [UserService],
)
class UserListComponent implements OnInit {
  final UserService _userService;

  UserListComponent(this._userService);

  List<User> users = [];

  @override
  void ngOnInit() {
    this._userService.list()
        .then((List<User> users) => this.users = users)
        .catchError((e) => print(e.toString()));
  }
}
