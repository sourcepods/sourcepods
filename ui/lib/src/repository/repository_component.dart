import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/src/gravatar_component.dart';
import 'package:gitpods/src/loading_component.dart';
import 'package:gitpods/src/repository/repository.dart';
import 'package:gitpods/src/repository/repository_service.dart';
import 'package:gitpods/src/repository/repository_tree.dart';

@Component(
  selector: 'gitpods-repository',
  templateUrl: 'repository_component.html',
  styleUrls: const ['repository_component.css'],
  directives: const [
    COMMON_DIRECTIVES,
    ROUTER_DIRECTIVES,
    LoadingComponent,
    GravatarComponent,
  ],
  providers: const [RepositoryService],
)
class RepositoryComponent implements OnInit {
  final RouteParams _routeParams;
  final RepositoryService _repositoryService;

  RepositoryComponent(this._routeParams, this._repositoryService);

  String ownerName;
  Repository repository;
  List<RepositoryTree> tree;

  @override
  void ngOnInit() {
    ownerName = this._routeParams.get('owner');
    String name = this._routeParams.get('name');

    this._repositoryService.get(ownerName, name).then((RepositoryPage page) {
      this.repository = page.repository;
    });
  }
}

class RepositoryPage {
  Repository repository;

  RepositoryPage(this.repository);
}
