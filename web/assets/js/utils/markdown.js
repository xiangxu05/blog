// Markdown 渲染工具 - markdown.js

class MarkdownRenderer {
  constructor() {
    this.renderer = null;
    this.initialized = false;
  }
  
  // 初始化Markdown渲染器
  init() {
    if (this.initialized) return;
    
    try {
      // 检查依赖库是否加载
      if (typeof marked === 'undefined') {
        console.error('Marked.js 未加载');
        return false;
      }
      
      // 创建自定义渲染器
      this.renderer = new marked.Renderer();
      
      // 自定义标题渲染（添加ID用于目录导航）
      this.renderer.heading = function(text, level) {
        const id = text.toLowerCase()
          .replace(/[^\w\u4e00-\u9fa5\s-]/g, '') // 移除特殊字符，保留中文、英文、数字、空格、连字符
          .replace(/\s+/g, '-') // 空格替换为连字符
          .replace(/-+/g, '-') // 多个连字符合并为一个
          .replace(/^-|-$/g, ''); // 移除首尾连字符
        
        return `<h${level} id="${id}" class="markdown-heading">
          <a href="#${id}" class="heading-anchor" aria-hidden="true">
            <i class="fas fa-link"></i>
          </a>
          ${text}
        </h${level}>`;
      };
      
      // 自定义代码块渲染
      this.renderer.code = function(code, language, escaped) {
        const langClass = language ? ` class="language-${language}"` : '';
        const langLabel = language ? `<div class="code-language">${language}</div>` : '';
        
        return `
          <div class="code-block-wrapper">
            ${langLabel}
            <button class="code-copy-btn" onclick="copyCodeBlock(this)" title="复制代码">
              <i class="fas fa-copy"></i>
            </button>
            <pre><code${langClass}>${escaped ? code : this.escape(code)}</code></pre>
          </div>
        `;
      };
      
      // 自定义内联代码渲染
      this.renderer.codespan = function(code) {
        return `<code class="inline-code">${code}</code>`;
      };
      
      // 自定义链接渲染（外链自动添加target="_blank"）
      this.renderer.link = function(href, title, text) {
        const isExternal = href.startsWith('http') && !href.includes(window.location.hostname);
        const target = isExternal ? ' target="_blank" rel="noopener noreferrer"' : '';
        const titleAttr = title ? ` title="${title}"` : '';
        const externalIcon = isExternal ? ' <i class="fas fa-external-link-alt fa-sm"></i>' : '';
        
        return `<a href="${href}"${titleAttr}${target}>${text}${externalIcon}</a>`;
      };
      
      // 自定义图片渲染
      this.renderer.image = function(href, title, text) {
        const titleAttr = title ? ` title="${title}"` : '';
        const altAttr = text ? ` alt="${text}"` : '';
        
        return `
          <figure class="image-figure">
            <img src="${href}"${altAttr}${titleAttr} class="markdown-image" loading="lazy">
            ${text ? `<figcaption>${text}</figcaption>` : ''}
          </figure>
        `;
      };
      
      // 自定义表格渲染
      this.renderer.table = function(header, body) {
        return `
          <div class="table-wrapper">
            <table class="table table-striped markdown-table">
              <thead>${header}</thead>
              <tbody>${body}</tbody>
            </table>
          </div>
        `;
      };
      
      // 自定义引用块渲染
      this.renderer.blockquote = function(quote) {
        return `
          <blockquote class="markdown-blockquote">
            <i class="fas fa-quote-left quote-icon"></i>
            ${quote}
          </blockquote>
        `;
      };
      
      // 自定义列表渲染
      this.renderer.list = function(body, ordered, start) {
        const type = ordered ? 'ol' : 'ul';
        const startAttr = ordered && start !== 1 ? ` start="${start}"` : '';
        return `<${type}${startAttr} class="markdown-list">${body}</${type}>`;
      };
      
      // 配置marked选项
      marked.setOptions({
        renderer: this.renderer,
        highlight: this.highlightCode.bind(this),
        breaks: true, // 支持GFM换行
        gfm: true, // 启用GitHub风格Markdown
        tables: true, // 支持表格
        sanitize: false, // 不删除HTML标签
        smartypants: true, // 智能标点符号
        xhtml: false
      });
      
      this.initialized = true;
      console.log('Markdown渲染器初始化完成');
      return true;
      
    } catch (error) {
      console.error('Markdown渲染器初始化失败:', error);
      return false;
    }
  }
  
  // 代码高亮处理
  highlightCode(code, language) {
    if (!window.hljs) {
      return this.escapeHtml(code);
    }
    
    try {
      if (language && window.hljs.getLanguage(language)) {
        return window.hljs.highlight(code, { language: language }).value;
      } else {
        return window.hljs.highlightAuto(code).value;
      }
    } catch (error) {
      console.warn('代码高亮失败:', error);
      return this.escapeHtml(code);
    }
  }
  
  // 渲染Markdown文本
  render(markdown) {
    if (!this.initialized) {
      if (!this.init()) {
        return `<p class="text-danger">Markdown渲染器初始化失败</p>`;
      }
    }
    
    try {
      const html = marked.parse(markdown);
      return html;
    } catch (error) {
      console.error('Markdown渲染失败:', error);
      return `
        <div class="alert alert-warning-custom">
          <i class="fas fa-exclamation-triangle me-2"></i>
          Markdown渲染失败，显示原始内容
        </div>
        <pre class="bg-light p-3 rounded">${this.escapeHtml(markdown)}</pre>
      `;
    }
  }
  
  // 渲染并处理数学公式
  async renderWithMath(markdown) {
    const html = this.render(markdown);
    
    // 如果有KaTeX库，渲染数学公式
    if (window.katex && window.renderMathInElement) {
      const tempDiv = document.createElement('div');
      tempDiv.innerHTML = html;
      
      try {
        window.renderMathInElement(tempDiv, {
          delimiters: [
            {left: '$$', right: '$$', display: true},
            {left: '$', right: '$', display: false},
            {left: '\\[', right: '\\]', display: true},
            {left: '\\(', right: '\\)', display: false}
          ],
          throwOnError: false
        });
        
        return tempDiv.innerHTML;
      } catch (error) {
        console.warn('数学公式渲染失败:', error);
        return html;
      }
    }
    
    return html;
  }
  
  // 生成目录
  generateTOC(markdown) {
    const headings = [];
    const renderer = new marked.Renderer();
    
    renderer.heading = function(text, level) {
      const id = text.toLowerCase()
        .replace(/[^\w\u4e00-\u9fa5\s-]/g, '')
        .replace(/\s+/g, '-')
        .replace(/-+/g, '-')
        .replace(/^-|-$/g, '');
      
      headings.push({
        level,
        text,
        id
      });
      
      return '';
    };
    
    // 临时设置渲染器来提取标题
    const originalRenderer = marked.defaults.renderer;
    marked.setOptions({ renderer });
    marked.parse(markdown);
    marked.setOptions({ renderer: originalRenderer });
    
    return headings;
  }
  
  // 计算阅读时间
  calculateReadingTime(markdown) {
    const wordsPerMinute = 500; // 每分钟阅读字数（中文）
    const wordCount = markdown.length;
    const readingTime = Math.ceil(wordCount / wordsPerMinute);
    
    return {
      wordCount,
      readingTime,
      readingTimeText: `${readingTime} 分钟阅读`
    };
  }
  
  // HTML转义
  escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }
}

// 全局工具函数

// 复制代码块
window.copyCodeBlock = function(button) {
  const codeBlock = button.parentElement.querySelector('code');
  const code = codeBlock.textContent;
  
  if (navigator.clipboard) {
    navigator.clipboard.writeText(code).then(() => {
      showCopySuccess(button);
    }).catch(() => {
      fallbackCopy(code, button);
    });
  } else {
    fallbackCopy(code, button);
  }
};

// 后备复制方法
function fallbackCopy(text, button) {
  const textArea = document.createElement('textarea');
  textArea.value = text;
  textArea.style.position = 'fixed';
  textArea.style.left = '-999999px';
  textArea.style.top = '-999999px';
  document.body.appendChild(textArea);
  textArea.focus();
  textArea.select();
  
  try {
    document.execCommand('copy');
    showCopySuccess(button);
  } catch (err) {
    console.error('复制失败:', err);
    window.blogApp?.showNotification('复制失败', 'error');
  }
  
  document.body.removeChild(textArea);
}

// 显示复制成功反馈
function showCopySuccess(button) {
  const icon = button.querySelector('i');
  const originalClass = icon.className;
  
  icon.className = 'fas fa-check';
  button.style.color = '#28a745';
  
  if (window.blogApp) {
    window.blogApp.showNotification('代码已复制到剪贴板', 'success');
  }
  
  setTimeout(() => {
    icon.className = originalClass;
    button.style.color = '';
  }, 2000);
}

// 切换代码高亮主题
window.toggleCodeHighlightTheme = function(theme) {
  const lightTheme = document.getElementById('highlight-theme-light');
  const darkTheme = document.getElementById('highlight-theme-dark');
  
  if (theme === 'dark') {
    lightTheme.disabled = true;
    darkTheme.disabled = false;
  } else {
    lightTheme.disabled = false;
    darkTheme.disabled = true;
  }
};

// 创建全局Markdown渲染器实例
window.markdownRenderer = new MarkdownRenderer();

// 监听主题变化，自动切换代码高亮主题
if (window.appState) {
  window.appState.subscribe('ui.theme', (theme) => {
    window.toggleCodeHighlightTheme(theme);
  });
}

// 初始化代码高亮主题
document.addEventListener('DOMContentLoaded', () => {
  const currentTheme = window.appState?.getState('ui.theme') || 'light';
  window.toggleCodeHighlightTheme(currentTheme);
});

// 导出（如果使用模块系统）
if (typeof module !== 'undefined' && module.exports) {
  module.exports = MarkdownRenderer;
}