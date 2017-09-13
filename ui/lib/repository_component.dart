import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/repository.dart';
import 'package:gitpods/repository_service.dart';

@Component(
  selector: 'gitpods-repository',
  templateUrl: 'repository_component.html',
  styleUrls: const['repository_component.css'],
  directives: const [COMMON_DIRECTIVES, ROUTER_DIRECTIVES],
  providers: const [RepositoryService],
)
class RepositoryComponent implements OnInit {
  final RouteParams _routeParams;
  final RepositoryService _repositoryService;

  RepositoryComponent(this._routeParams, this._repositoryService);

  String ownerName;
  Repository repository;

  @override
  ngOnInit() {
    ownerName = this._routeParams.get('owner');
    String name = this._routeParams.get('name');

    this._repositoryService.get(ownerName, name)
        .then((Repository repository) => this.repository = repository);
  }
}
