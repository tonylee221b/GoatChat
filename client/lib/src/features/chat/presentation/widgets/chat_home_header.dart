import 'package:flutter/material.dart';

import '../../../../theme/app_theme.dart';
import 'circle_icon_button.dart';

class ChatHomeHeader extends StatelessWidget {
  const ChatHomeHeader({super.key});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        const Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Hello,',
                style: TextStyle(color: AppColors.mutedText, fontSize: 14),
              ),
              SizedBox(height: 4),
              Text(
                'GoatChat',
                style: TextStyle(
                  color: AppColors.textPrimary,
                  fontSize: 30,
                  fontWeight: FontWeight.w800,
                ),
              ),
            ],
          ),
        ),
        CircleIconButton(icon: Icons.search_rounded, onPressed: () {}),
        const SizedBox(width: 10),
        CircleIconButton(icon: Icons.more_vert_rounded, onPressed: () {}),
      ],
    );
  }
}
