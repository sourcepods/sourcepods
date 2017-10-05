import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/gravatar_component.dart';
import 'package:gitpods/repository.dart';
import 'package:gitpods/repository_service.dart';
import 'package:gitpods/repository_tree.dart';

@Component(
  selector: 'gitpods-repository',
  templateUrl: 'repository_component.html',
  styleUrls: const['repository_component.css'],
  directives: const [COMMON_DIRECTIVES, ROUTER_DIRECTIVES, Gravatar],
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
  ngOnInit() {
    ownerName = this._routeParams.get('owner');
    String name = this._routeParams.get('name');

    this._repositoryService.get(ownerName, name).then((RepositoryPage page) {
      this.repository = page.repository;

      var tree = page.tree;
      tree.sort((a, b) {
        if (a.type == b.type) {
          return a.file.compareTo(b.file);
        } else {
          return b.type.compareTo(a.type);
        }
      });
      this.tree = tree;
    });
  }
}

class RepositoryPage {
  Repository repository;
  List<RepositoryTree> tree;

  RepositoryPage(this.repository, this.tree);
}