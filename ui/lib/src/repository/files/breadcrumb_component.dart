import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';

@Component(
  selector: 'files-breadcrumb',
  templateUrl: 'breadcrumb_component.html',
  styleUrls: ['breadcrumb_component.css'],
  directives: [
    coreDirectives,
    routerDirectives,
  ],
)
class BreadcrumbComponent implements OnChanges {
  @Input()
  String ownerName;
  @Input()
  String repositoryName;

  @Input()
  String path;

  List<String> elements;

  @override
  void ngOnChanges(Map<String, SimpleChange> changes) {
    this.elements = this.path.split('/');
  }

  String changeRoot() {
    return './$ownerName/$repositoryName';
  }

  String changeDirectory(int i) {
    String currentBranch = 'master'; // TODO
    String newPath = path.split('/').getRange(0, i+1).join('/');
    return './$ownerName/$repositoryName/tree/$currentBranch/$newPath';
  }
}
