import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:sourcepods/footer_component.dart';
import 'package:sourcepods/header_component.dart';
import 'package:sourcepods/routes.dart';

@Component(
  selector: 'sourcepods-app',
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
