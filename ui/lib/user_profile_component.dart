import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/gravatar_component.dart';
import 'package:gitpods/mailto_pipe.dart';
import 'package:gitpods/user.dart';
import 'package:gitpods/user_service.dart';

@Component(
  selector: 'gitpods-user-profile',
  templateUrl: 'user_profile_component.html',
  styleUrls: const ['user_profile_component.css'],
  providers: const [UserService],
  directives: const [COMMON_DIRECTIVES, Gravatar],
  pipes: const [DatePipe, MailtoPipe],
)
class UserProfileComponent implements OnInit {
  final RouteParams _routeParams;
  final UserService _userService;

  UserProfileComponent(this._routeParams, this._userService);

  User user;

  @override
  ngOnInit() {
    String username = this._routeParams.get('username');
    this._userService.profile(username)
        .then((user) => this.user = user);
  }
}
