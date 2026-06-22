import 'package:flutter/material.dart';

import '../../../theme/app_theme.dart';
import 'widgets/chat_filter_tabs.dart';
import 'widgets/chat_home_header.dart';

class ChatHomePage extends StatelessWidget {
  const ChatHomePage({required this.nickname, super.key});

  final String nickname;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.fromLTRB(20, 16, 20, 20),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              ChatHomeHeader(nickname: nickname),
              const SizedBox(height: 24),
              const ChatFilterTabs(),
              const SizedBox(height: 20),
              const Expanded(child: SizedBox()),
            ],
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {},
        backgroundColor: AppColors.ink,
        foregroundColor: Colors.white,
        shape: const CircleBorder(),
        child: const Icon(Icons.chat_bubble_outline_rounded),
      ),
    );
  }
}
