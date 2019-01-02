import 'package:angular/angular.dart';
import 'package:angular_forms/angular_forms.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/src/gravatar_component.dart';
import 'package:gitpods/src/loading_component.dart';
import 'package:gitpods/src/mailto_pipe.dart';
import 'package:gitpods/src/repository/repository.dart';
import 'package:gitpods/src/user/user.dart';
import 'package:gitpods/src/user/user_service.dart';

@Component(
  selector: 'gitpods-user-profile',
  templateUrl: 'user_profile_component.html',
  styleUrls: const ['user_profile_component.css'],
  providers: const [UserService],
  directives: const [
    coreDirectives,
    routerDirectives,
    formDirectives,
    LoadingComponent,
    GravatarComponent,
  ],
  pipes: const [DatePipe, MailtoPipe, FilteredReposPipe],
)
class UserProfileComponent implements OnActivate {
  UserProfileComponent(this._userService);

  final UserService _userService;

  User user;
  List<Repository> repositories;
  String repoQuery = '';

  @override
  void onActivate(RouterState previous, RouterState current) {
    String username = current.parameters['username'];
    this._userService.profile(username).then((UserProfile profile) {
      this.user = profile.user;
      this.repositories = profile.repositories;
    });
  }
}

@Pipe('filteredRepos')
class FilteredReposPipe extends PipeTransform {
  List<Repository> transform(List<Repository> repos, String query) {
    repos = repos.where((repo) => repo.name.contains(query)).toList();
    repos.sort((a, b) => a.name.compareTo(b.name));
    return repos;
  }
}
