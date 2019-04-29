import {Injectable} from '@nestjs/common';
import {Category} from './category.entity';
import {Repository} from 'typeorm';
import {InjectRepository} from '@nestjs/typeorm';

@Injectable()
export class CategoryService {

  constructor(
    @InjectRepository(Category)
    private readonly categoryRepository: Repository<Category>,
  ) {
  }

  public async findAllCategories(): Promise<Category[]> {
    return this.categoryRepository.find();
  }

  public async createCategory(category: Category) {
    return this.categoryRepository.insert(category);
  }
}
