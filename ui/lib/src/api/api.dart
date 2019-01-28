library swagger.api;

import 'dart:async';
import 'dart:convert';
import 'package:http/browser_client.dart';
import 'package:http/http.dart';

part 'api_client.dart';
part 'api_helper.dart';
part 'api_exception.dart';
part 'auth/authentication.dart';
part 'auth/api_key_auth.dart';
part 'auth/oauth.dart';
part 'auth/http_basic_auth.dart';

part 'api/repositories_api.dart';
part 'api/users_api.dart';

part 'model/branch.dart';
part 'model/error.dart';
part 'model/new_repository.dart';
part 'model/repository.dart';
part 'model/tree_entry.dart';
part 'model/updated_user.dart';
part 'model/user.dart';
part 'model/validation_error.dart';
part 'model/validation_error_errors.dart';


ApiClient defaultApiClient = new ApiClient();
