import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/routes.dart' as global;
import 'package:gitpods/src/gravatar_component.dart';
import 'package:gitpods/src/loading_component.dart';
import 'package:gitpods/src/repository/repository.dart';
import 'package:gitpods/src/repository/repository_service.dart';
import 'package:gitpods/src/repository/routes.dart';

@Component(
  selector: 'gitpods-repository',
  templateUrl: 'repository_component.html',
  styleUrls: const ['repository_component.css'],
  directives: [
    coreDirectives,
    routerDirectives,
    LoadingComponent,
    GravatarComponent,
  ],
  providers: const [
    ClassProvider(RepositoryService),
    ClassProvider(Routes),
  ],
  exports: [Routes],
)
class RepositoryComponent implements OnActivate {
  RepositoryComponent(this._repositoryService);

  final RepositoryService _repositoryService;

  String ownerName;
  Repository repository;

  @override
  void onActivate(RouterState previous, RouterState current) {
    ownerName = current.parameters['owner'];
    String name = current.parameters['name'];

    _repositoryService.get(ownerName, name).then((r) => this.repository = r);
  }

  String userProfileUrl() =>
      global.RoutePaths.userProfile.toUrl(parameters: {'username': ownerName});

  Map<String, String> _parameters() => {
        'owner': this.ownerName,
        'name': this.repository.name,
      };

  String repositoryUrl() =>
      global.RoutePaths.repository.toUrl(parameters: _parameters());

  String filesUrl() => RoutePaths.files.toUrl(parameters: _parameters());

  String commitsUrl() => RoutePaths.commits.toUrl(parameters: _parameters());

  String settingsUrl() => RoutePaths.settings.toUrl(parameters: _parameters());
}
