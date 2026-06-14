import 'package:flutter/material.dart';

import 'features/chat/presentation/chat_home_page.dart';
import 'theme/app_theme.dart';

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'GoatChat',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.light,
      home: const ChatHomePage(),
    );
  }
}
