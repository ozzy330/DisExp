#ifndef EXTRA_H
#define EXTRA_H

#include "qoi.h"

/* Encode raw RGB or RGBA pixels into a QOI image in memory.

The function either returns NULL on failure (invalid parameters or malloc
failed) or a pointer to the encoded data on success. On success the out_len
is set to the size in bytes of the encoded data.

The returned qoi data should be free()d after use. */

void *qoi_encode_run(const void *data, const qoi_desc *desc, int *out_len);
void *qoi_encode_diff_luma(const void *data, const qoi_desc *desc,
                           int *out_len);

#ifdef QOI_IMPLEMENTATION
void *qoi_encode_diff_luma(const void *data, const qoi_desc *desc,
                           int *out_len) {
  int i, max_size, p, run;
  int px_len, px_end, px_pos, channels;
  unsigned char *bytes;
  const unsigned char *pixels;
  qoi_rgba_t index[64];
  qoi_rgba_t px, px_prev;

  if (data == NULL || out_len == NULL || desc == NULL || desc->width == 0 ||
      desc->height == 0 || desc->channels < 3 || desc->channels > 4 ||
      desc->colorspace > 1 || desc->height >= QOI_PIXELS_MAX / desc->width) {
    return NULL;
  }

  max_size = desc->width * desc->height * (desc->channels + 1) +
             QOI_HEADER_SIZE + sizeof(qoi_padding);

  p = 0;
  bytes = (unsigned char *)QOI_MALLOC(max_size);
  if (!bytes) {
    return NULL;
  }

  qoi_write_32(bytes, &p, QOI_MAGIC);
  qoi_write_32(bytes, &p, desc->width);
  qoi_write_32(bytes, &p, desc->height);
  bytes[p++] = desc->channels;
  bytes[p++] = desc->colorspace;

  pixels = (const unsigned char *)data;

  QOI_ZEROARR(index);

  run = 0;
  px_prev.rgba.r = 0;
  px_prev.rgba.g = 0;
  px_prev.rgba.b = 0;
  px_prev.rgba.a = 255;
  px = px_prev;

  px_len = desc->width * desc->height * desc->channels;
  px_end = px_len - desc->channels;
  channels = desc->channels;

  for (px_pos = 0; px_pos < px_len; px_pos += channels) {
    px.rgba.r = pixels[px_pos + 0];
    px.rgba.g = pixels[px_pos + 1];
    px.rgba.b = pixels[px_pos + 2];

    if (channels == 4) {
      px.rgba.a = pixels[px_pos + 3];
    }

    /*WARN: This code will not do RUN or Index encoding*/
    /*if (px.v == px_prev.v) {*/
    /*	run++;*/
    /*	if (run == 62 || px_pos == px_end) {*/
    /*		bytes[p++] = QOI_OP_RUN | (run - 1);*/
    /*		run = 0;*/
    /*	}*/
    /*}*/
    /*else {*/
    /*int index_pos;*/

    /*WARN: This code will not do RUN or Index encoding*/
    /*if (run > 0) {*/
    /*  bytes[p++] = QOI_OP_RUN | (run - 1);*/
    /*  run = 0;*/
    /*}*/
    /*index_pos = QOI_COLOR_HASH(px) % 64;*/
    /*if (index[index_pos].v == px.v) {*/
    /*  bytes[p++] = QOI_OP_INDEX | index_pos;*/
    /*} else {*/
    /*  index[index_pos] = px;*/

    if (px.rgba.a == px_prev.rgba.a) {
      signed char vr = px.rgba.r - px_prev.rgba.r;
      signed char vg = px.rgba.g - px_prev.rgba.g;
      signed char vb = px.rgba.b - px_prev.rgba.b;

      signed char vg_r = vr - vg;
      signed char vg_b = vb - vg;

      if (vr > -3 && vr < 2 && vg > -3 && vg < 2 && vb > -3 && vb < 2) {
        bytes[p++] = QOI_OP_DIFF | (vr + 2) << 4 | (vg + 2) << 2 | (vb + 2);
      } else if (vg_r > -9 && vg_r < 8 && vg > -33 && vg < 32 && vg_b > -9 &&
                 vg_b < 8) {
        bytes[p++] = QOI_OP_LUMA | (vg + 32);
        bytes[p++] = (vg_r + 8) << 4 | (vg_b + 8);
      } else {
        bytes[p++] = QOI_OP_RGB;
        bytes[p++] = px.rgba.r;
        bytes[p++] = px.rgba.g;
        bytes[p++] = px.rgba.b;
      }
    } else {
      bytes[p++] = QOI_OP_RGBA;
      bytes[p++] = px.rgba.r;
      bytes[p++] = px.rgba.g;
      bytes[p++] = px.rgba.b;
      bytes[p++] = px.rgba.a;
    }
    /*WARN: This code will not do RUN or Index encoding*/
    /*}*/
    /*}*/
    px_prev = px;
  }

  for (i = 0; i < (int)sizeof(qoi_padding); i++) {
    bytes[p++] = qoi_padding[i];
  }

  *out_len = p;
  return bytes;
}

void *qoi_encode_run(const void *data, const qoi_desc *desc, int *out_len) {
  int i, max_size, p, run;
  int px_len, px_end, px_pos, channels;
  unsigned char *bytes;
  const unsigned char *pixels;
  qoi_rgba_t index[64];
  qoi_rgba_t px, px_prev;

  if (data == NULL || out_len == NULL || desc == NULL || desc->width == 0 ||
      desc->height == 0 || desc->channels < 3 || desc->channels > 4 ||
      desc->colorspace > 1 || desc->height >= QOI_PIXELS_MAX / desc->width) {
    return NULL;
  }

  max_size = desc->width * desc->height * (desc->channels + 1) +
             QOI_HEADER_SIZE + sizeof(qoi_padding);

  p = 0;
  bytes = (unsigned char *)QOI_MALLOC(max_size);
  if (!bytes) {
    return NULL;
  }

  qoi_write_32(bytes, &p, QOI_MAGIC);
  qoi_write_32(bytes, &p, desc->width);
  qoi_write_32(bytes, &p, desc->height);
  bytes[p++] = desc->channels;
  bytes[p++] = desc->colorspace;

  pixels = (const unsigned char *)data;

  QOI_ZEROARR(index);

  run = 0;
  px_prev.rgba.r = 0;
  px_prev.rgba.g = 0;
  px_prev.rgba.b = 0;
  px_prev.rgba.a = 255;
  px = px_prev;

  px_len = desc->width * desc->height * desc->channels;
  px_end = px_len - desc->channels;
  channels = desc->channels;

  for (px_pos = 0; px_pos < px_len; px_pos += channels) {
    px.rgba.r = pixels[px_pos + 0];
    px.rgba.g = pixels[px_pos + 1];
    px.rgba.b = pixels[px_pos + 2];

    if (channels == 4) {
      px.rgba.a = pixels[px_pos + 3];
    }

    if (px.v == px_prev.v) {
      run++;
      if (run == 62 || px_pos == px_end) {
        bytes[p++] = QOI_OP_RUN | (run - 1);
        run = 0;
      }
    } else {
      /*WARN: This code will not do Diff/Luma or Index encoding*/
      /*int index_pos;*/

      if (run > 0) {
        bytes[p++] = QOI_OP_RUN | (run - 1);
        run = 0;
      }

      /*WARN: This code will not do Diff/Luma or Index encoding*/
      /*index_pos = QOI_COLOR_HASH(px) % 64;*/

      /*WARN: This code will not do Diff/Luma or Index encoding*/
      /*if (index[index_pos].v == px.v) {*/
      /*  bytes[p++] = QOI_OP_INDEX | index_pos;*/
      /*} else {*/
      /*  index[index_pos] = px;*/

      /*WARN: This code will not do Diff/Luma or Index encoding*/
      /*if (px.rgba.a == px_prev.rgba.a) {*/
      /*  signed char vr = px.rgba.r - px_prev.rgba.r;*/
      /*  signed char vg = px.rgba.g - px_prev.rgba.g;*/
      /*  signed char vb = px.rgba.b - px_prev.rgba.b;*/
      /**/
      /*  signed char vg_r = vr - vg;*/
      /*  signed char vg_b = vb - vg;*/
      /**/
      /*  if (vr > -3 && vr < 2 && vg > -3 && vg < 2 && vb > -3 && vb < 2) {*/
      /*    bytes[p++] = QOI_OP_DIFF | (vr + 2) << 4 | (vg + 2) << 2 | (vb +
       * 2);*/
      /*  } else if (vg_r > -9 && vg_r < 8 && vg > -33 && vg < 32 && vg_b > -9
       * &&*/
      /*             vg_b < 8) {*/
      /*    bytes[p++] = QOI_OP_LUMA | (vg + 32);*/
      /*    bytes[p++] = (vg_r + 8) << 4 | (vg_b + 8);*/
      /*  } else {*/
      /*    bytes[p++] = QOI_OP_RGB;*/
      /*    bytes[p++] = px.rgba.r;*/
      /*    bytes[p++] = px.rgba.g;*/
      /*    bytes[p++] = px.rgba.b;*/
      /*  }*/
      /*} else {*/
      bytes[p++] = QOI_OP_RGBA;
      bytes[p++] = px.rgba.r;
      bytes[p++] = px.rgba.g;
      bytes[p++] = px.rgba.b;
      bytes[p++] = px.rgba.a;
      /*WARN: This code will not do Diff/Luma or Index encoding*/
      /*}*/
      /*}*/
    }
    px_prev = px;
  }

  for (i = 0; i < (int)sizeof(qoi_padding); i++) {
    bytes[p++] = qoi_padding[i];
  }

  *out_len = p;
  return bytes;
}

void *qoi_encode_index(const void *data, const qoi_desc *desc, int *out_len) {
  int i, max_size, p, run;
  int px_len, px_end, px_pos, channels;
  unsigned char *bytes;
  const unsigned char *pixels;
  qoi_rgba_t index[64];
  qoi_rgba_t px, px_prev;

  if (data == NULL || out_len == NULL || desc == NULL || desc->width == 0 ||
      desc->height == 0 || desc->channels < 3 || desc->channels > 4 ||
      desc->colorspace > 1 || desc->height >= QOI_PIXELS_MAX / desc->width) {
    return NULL;
  }

  max_size = desc->width * desc->height * (desc->channels + 1) +
             QOI_HEADER_SIZE + sizeof(qoi_padding);

  p = 0;
  bytes = (unsigned char *)QOI_MALLOC(max_size);
  if (!bytes) {
    return NULL;
  }

  qoi_write_32(bytes, &p, QOI_MAGIC);
  qoi_write_32(bytes, &p, desc->width);
  qoi_write_32(bytes, &p, desc->height);
  bytes[p++] = desc->channels;
  bytes[p++] = desc->colorspace;

  pixels = (const unsigned char *)data;

  QOI_ZEROARR(index);

  run = 0;
  px_prev.rgba.r = 0;
  px_prev.rgba.g = 0;
  px_prev.rgba.b = 0;
  px_prev.rgba.a = 255;
  px = px_prev;

  px_len = desc->width * desc->height * desc->channels;
  px_end = px_len - desc->channels;
  channels = desc->channels;

  for (px_pos = 0; px_pos < px_len; px_pos += channels) {
    px.rgba.r = pixels[px_pos + 0];
    px.rgba.g = pixels[px_pos + 1];
    px.rgba.b = pixels[px_pos + 2];

    if (channels == 4) {
      px.rgba.a = pixels[px_pos + 3];
    }

    int index_pos;

    index_pos = QOI_COLOR_HASH(px) % 64;

    if (index[index_pos].v == px.v) {
      bytes[p++] = QOI_OP_INDEX | index_pos;
    } else {
      index[index_pos] = px;
      bytes[p++] = QOI_OP_RGBA;
      bytes[p++] = px.rgba.r;
      bytes[p++] = px.rgba.g;
      bytes[p++] = px.rgba.b;
      bytes[p++] = px.rgba.a;
    }
    px_prev = px;
  }

  for (i = 0; i < (int)sizeof(qoi_padding); i++) {
    bytes[p++] = qoi_padding[i];
  }

  *out_len = p;
  return bytes;
}

#endif /* QOI_IMPLEMENTATION */
#endif /* EXTRA_H */
