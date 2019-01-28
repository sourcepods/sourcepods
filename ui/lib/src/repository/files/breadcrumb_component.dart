import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';

@Component(
  selector: 'files-breadcrumb',
  templateUrl: 'breadcrumb_component.html',
  directives: [
    coreDirectives,
    routerDirectives,
  ],
)
class BreadcrumbComponent implements OnChanges{
  @Input()
  String path;

  List<String> elements;

  @override
  void ngOnChanges(Map<String, SimpleChange> changes) {
    this.elements = this.path.split('/');
  }

  String changeDirectory(element) {
    return '';
  }
}
