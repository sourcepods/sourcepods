import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/footer_component.dart';
import 'package:gitpods/header_component.dart';
import 'package:gitpods/routes.dart';

@Component(
  selector: 'gitpods-app',
  templateUrl: 'app_component.html',
  styleUrls: const ['app_component.css'],
  directives: const [
    coreDirectives,
    routerDirectives,
    HeaderComponent,
    FooterComponent,
  ],
  providers: const [
    routerProviders,
    ClassProvider(Routes),
  ],
  exports: [
    RoutePaths,
    Routes,
  ],
)
class AppComponent {}
